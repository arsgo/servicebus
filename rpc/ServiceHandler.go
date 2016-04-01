package rpc

type serviceRequestSnap struct {
	total   int
	current int
	success int
	failed  int
}

type serviceCallBack interface {
	Request(string, string) (string, error)
	Send(string, string, []byte) (string, error)
}

type serviceHandler struct {
	services map[string]*serviceRequestSnap
	callback serviceCallBack
}

func NewServiceHandler(callback serviceCallBack) *serviceHandler {
	return &serviceHandler{services: make(map[string]*serviceRequestSnap), callback: callback}
}

func (s *serviceHandler) Request(name string, input string) (r string, err error) {
	if _, ok := s.services[name]; !ok {
		s.services[name] = &serviceRequestSnap{}
	}
	s.services[name].total++
	s.services[name].current++
	defer func() {
		s.services[name].current--
		if err == nil {
			s.services[name].success++
		} else {
			s.services[name].failed++
		}
	}()
	return s.callback.Request(name, input)
}

func (s *serviceHandler) Send(name string, input string, data []byte) (r string, err error) {
	if _, ok := s.services[name]; !ok {
		s.services[name] = &serviceRequestSnap{}
	}
	s.services[name].total++
	s.services[name].current++
	defer func() {
		s.services[name].current--
		if err == nil {
			s.services[name].success++
		} else {
			s.services[name].failed++
		}
	}()
	return s.callback.Send(name, input,data)
}
