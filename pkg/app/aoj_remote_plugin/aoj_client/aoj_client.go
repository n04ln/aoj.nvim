package aoj_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	ErrSomethingWrong = errors.New("something wrong")
)

type AojClient struct {
	cookies []*http.Cookie
	mux     sync.Mutex

	c *resty.Client
}

func NewAojClient() *AojClient {
	c := resty.New().
		SetBaseURL("https://judgeapi.u-aizu.ac.jp").
		SetHeader("Content-Type", "application/json")
	return &AojClient{
		c: c,
	}
}

func (a *AojClient) Session(ctx context.Context, id, password string) error {
	req := struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}{
		id, password,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := a.c.R().SetContext(ctx).
		SetBody(body).
		Post("/session")
	if err != nil {
		return err
	}

	a.mux.Lock()
	defer a.mux.Unlock()
	a.cookies = res.Cookies()

	return err
}

func (a *AojClient) GetDescription(ctx context.Context, problemID string) (*Description, error) {
	path := fmt.Sprintf("/resources/descriptions/en/%s", problemID)
	var description Description

	res, err := a.c.R().SetContext(ctx).
		SetResult(&description).
		Get(path)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, ErrSomethingWrong
	}

	return &description, nil
}

func (a *AojClient) Submit(ctx context.Context, problemID, language, sourceCode string) (string, error) {
	req := struct {
		ProblemID  string `json:"problemId"`
		Language   string `json:"language"`
		SourceCode string `json:"sourceCode"`
	}{
		problemID, language, sourceCode,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	result := SubmitResponse{}
	res, err := a.c.R().SetContext(ctx).
		SetBody(body).
		SetResult(&result).
		SetCookies(a.cookies).
		Post("/submissions")
	if err != nil {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", ErrSomethingWrong
	}

	return result.Token, nil
}

func (a *AojClient) Status(ctx context.Context, token, problemID string) (*Status, error) {
	result := []RecentSubmission{}
	res, err := a.c.R().SetContext(ctx).
		SetResult(&result).
		SetCookies(a.cookies).
		Get("/submission_records/recent")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, ErrSomethingWrong
	}

	var submission RecentSubmission
	for _, r := range result {
		if r.Token == token {
			submission = r
		}
	}

	const maxRetryCount = 100
	path := fmt.Sprintf("/verdicts/%d", submission.JudgeId)
	for i := 0; i < maxRetryCount; i++ {
		time.Sleep(100 * time.Millisecond)
		result := Status{}
		res, err := a.c.R().SetContext(ctx).
			SetResult(&result).
			Get(path)
		if err != nil {
			continue
		}
		if res.StatusCode() != http.StatusOK {
			continue
		}

		return &result, nil
	}

	return nil, errors.New("failed get verdicts")
}
