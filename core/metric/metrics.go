package metric

type Metrics struct {
	Results map[string][]*Result
}

func (m *Metrics) ConsumeResult(result *Result) {
	if m.Results == nil {
		m.Results = make(map[string][]*Result, 0)
	}

	if _, ok := m.Results[result.Name]; !ok {
		m.Results[result.Name] = make([]*Result, 0)
	}

	m.Results[result.Name] = append(m.Results[result.Name], result)
}
