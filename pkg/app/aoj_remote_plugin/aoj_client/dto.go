package aoj_client

import (
	"fmt"
	"strings"
	"time"
)

type Commentary struct {
	Filter  []string `json:"filter"`
	Pattern string   `json:"pattern"`
	Type    string   `json:"type"`
}

type Description struct {
	Commentaries    []*Commentary `json:"commentaries"`
	CreatedAt       int64         `json:"created_at"`
	HTML            string        `json:"html"`
	IsSolved        bool          `json:"isSolved"`
	Language        string        `json:"language"`
	MemoryLimit     int           `json:"memory_limit"`
	ProblemId       string        `json:"problem_id"`
	Recommendations int           `json:"recommendations"`
	Score           float64       `json:"score"`
	ServerTime      int           `json:"server_time"`
	SolvedUser      int           `json:"solvedUser"`
	SuccessRate     float64       `json:"successRate"`
	TimeLimit       int           `json:"time_limit"`
}

func (d *Description) String() string {
	tmpl := "ProblemId: %v\nIsSolved?: %v\nCreatedAt: %v\nTimeLimit: %v sec\n\nplease see below:\n http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=%v"

	var problemId string
	problemId = d.ProblemId

	var isSolvedMsg string
	if d.IsSolved {
		isSolvedMsg = "yes"
	} else {
		isSolvedMsg = "no"
	}

	t := time.Unix(d.CreatedAt/1000, 0)
	var createdAt string
	createdAt = t.Format("2006-01-02")

	var timeLimit int
	timeLimit = d.TimeLimit

	return fmt.Sprintf(tmpl, problemId, isSolvedMsg, createdAt, timeLimit, problemId)
}

type CaseVerdicts struct {
	CpuTime int64  `json:"cpuTime"`
	Memory  int64  `json:"memory"`
	Serial  int    `json:"serial"`
	Label   string `json:"label"`
	Status  string `json:"status"`
}

type Status struct {
	CaseVerdicts []CaseVerdicts `json:"caseVerdicts"`
	CompileError string         `json:"compileError"`
	RuntimeError string         `json:"runtimeError"`
	UserOutput   string         `json:"userOutput"`
	ProblemId    string
	Time         time.Time
}

func (s *Status) String() string {
	isAllCasesAC := true
	caseVerdictTemp := `testcase: %v, Memory: %vkB, CpuTime: %vs, Status: %v`
	messages := make([]string, 0, len(s.CaseVerdicts)+3)
	for _, cv := range s.CaseVerdicts {
		messages = append(messages, fmt.Sprintf(caseVerdictTemp, cv.Label, cv.Memory, cv.CpuTime, cv.Status))
		if cv.Status != "AC" {
			isAllCasesAC = false
		}
	}

	if s.CompileError != "" {
		messages = append(messages, fmt.Sprintf("CompileError: %v", s.CompileError))
	}

	if s.RuntimeError != "" {
		messages = append(messages, fmt.Sprintf("RuntimeError: %v", s.RuntimeError))
	}

	if s.UserOutput != "" {
		messages = append(messages, fmt.Sprintf("UserOutput: %v", s.UserOutput))
	}

	var comment string
	if isAllCasesAC {
		comment = "\n✅ CONGRATULATION!"
	} else {
		comment = "\n❌ NOT BE ACCEPTED!"
	}

	return "Submission result:\n" + strings.Join(messages, "\n") + comment
}

type SubmitResponse struct {
	Token string `json:"token"`
}

type RecentSubmission struct {
	Accuracy       string `json:"accuracy"`
	CodeSize       int    `json:"codeSize"`
	CpuTime        int    `json:"cpuTime"`
	JudgeDate      int64  `json:"judgeDate"`
	JudgeId        int64  `json:"judgeId"`
	JudgeType      int    `json:"judgeType"`
	Language       string `json:"language"`
	Memory         int    `json:"memory"`
	ProblemId      string `json:"problemId"`
	ProblemTitle   string `json:"problemTitle"`
	Status         int    `json:"status"`
	SubmissionDate int64  `json:"submissionDate"`
	Token          string `json:"token"`
	UserId         string `json:"userId"`
}
