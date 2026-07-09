package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	pncp "github.com/danyele/podp/internal/shared/clients/pncp"
)

type jobState struct {
	eventChan   chan pncp.EventoAnalise
	results     []*pncp.AnaliseResultado
	paginasErro []int
	mu          sync.RWMutex
	done        chan struct{}
}

type JobManager struct {
	mu   sync.Mutex
	jobs map[string]*jobState
}

func NovoJobManager() *JobManager {
	return &JobManager{jobs: make(map[string]*jobState)}
}

func (jm *JobManager) CriarJob() (string, chan pncp.EventoAnalise, *jobState) {
	jobID := uuid.New().String()
	eventChan := make(chan pncp.EventoAnalise, 200)
	js := &jobState{eventChan: eventChan, done: make(chan struct{})}
	jm.mu.Lock()
	jm.jobs[jobID] = js
	jm.mu.Unlock()
	return jobID, eventChan, js
}

func (jm *JobManager) FinalizarJob(jobID string, results []*pncp.AnaliseResultado, paginasErro []int) {
	jm.mu.Lock()
	if job, ok := jm.jobs[jobID]; ok {
		job.mu.Lock()
		job.results = results
		job.paginasErro = paginasErro
		job.mu.Unlock()
		close(job.done)
	}
	jm.mu.Unlock()
}

func (jm *JobManager) GetJobChan(jobID string) (<-chan pncp.EventoAnalise, bool) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	job, exists := jm.jobs[jobID]
	if !exists {
		return nil, false
	}
	return job.eventChan, true
}

func (jm *JobManager) BuscarResultadosBatch(c *gin.Context) {
	jobID := c.Param("jobId")

	jm.mu.Lock()
	job, exists := jm.jobs[jobID]
	jm.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"erro": "job nao encontrado"})
		return
	}

	<-job.done

	job.mu.RLock()
	results := job.results
	paginasErro := job.paginasErro
	job.mu.RUnlock()

	status := "processing"
	if results != nil {
		status = "completed"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      status,
		"results":     results,
		"paginasErro": paginasErro,
	})
}
