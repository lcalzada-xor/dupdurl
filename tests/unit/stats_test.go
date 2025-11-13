package unit

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/lcalzada-xor/dupdurl/pkg/stats"
)

func TestStatisticsBasic(t *testing.T) {
	st := stats.NewStatistics()

	st.TotalProcessed = 100
	st.UniqueURLs = 20
	st.Duplicates = 75
	st.ParseErrors = 3
	st.Filtered = 2

	if st.TotalProcessed != 100 {
		t.Errorf("TotalProcessed = %d; want 100", st.TotalProcessed)
	}
}

func TestProcessingTime(t *testing.T) {
	st := stats.NewStatistics()

	// Simulate some processing time
	time.Sleep(10 * time.Millisecond)
	st.Finish()

	duration := st.ProcessingTime()
	if duration < 10*time.Millisecond {
		t.Errorf("ProcessingTime() = %v; want >= 10ms", duration)
	}
}

func TestPrintStatistics(t *testing.T) {
	st := stats.NewStatistics()
	st.TotalProcessed = 100
	st.UniqueURLs = 20
	st.Duplicates = 80

	var buf bytes.Buffer
	st.Print(&buf)

	output := buf.String()
	if !strings.Contains(output, "100") {
		t.Error("Output should contain total processed count")
	}
	if !strings.Contains(output, "20") {
		t.Error("Output should contain unique count")
	}
	if !strings.Contains(output, "80") {
		t.Error("Output should contain duplicates count")
	}
}

func TestRecordDomain(t *testing.T) {
	st := stats.NewStatistics()

	st.RecordDomain("example.com")
	st.RecordDomain("example.com")
	st.RecordDomain("test.com")

	if st.TopDomains["example.com"] != 2 {
		t.Errorf("TopDomains[example.com] = %d; want 2", st.TopDomains["example.com"])
	}
	if st.TopDomains["test.com"] != 1 {
		t.Errorf("TopDomains[test.com] = %d; want 1", st.TopDomains["test.com"])
	}
}

func TestRecordParam(t *testing.T) {
	st := stats.NewStatistics()

	st.RecordParam("id")
	st.RecordParam("sort")
	st.RecordParam("id")

	if st.ParamFrequency["id"] != 2 {
		t.Errorf("ParamFrequency[id] = %d; want 2", st.ParamFrequency["id"])
	}
	if st.ParamFrequency["sort"] != 1 {
		t.Errorf("ParamFrequency[sort] = %d; want 1", st.ParamFrequency["sort"])
	}
}

func TestAvgQueryParams(t *testing.T) {
	st := stats.NewStatistics()
	st.UniqueURLs = 10

	for i := 0; i < 30; i++ {
		st.RecordParam("param")
	}

	avg := st.AvgQueryParams()
	if avg != 3.0 {
		t.Errorf("AvgQueryParams() = %f; want 3.0", avg)
	}
}

func TestToJSON(t *testing.T) {
	st := stats.NewStatistics()
	st.TotalProcessed = 100
	st.UniqueURLs = 20

	jsonData := st.ToJSON()

	if jsonData["total_processed"] != 100 {
		t.Errorf("JSON total_processed = %v; want 100", jsonData["total_processed"])
	}
	if jsonData["unique_urls"] != 20 {
		t.Errorf("JSON unique_urls = %v; want 20", jsonData["unique_urls"])
	}
}
