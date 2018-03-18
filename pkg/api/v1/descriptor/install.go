package descriptor

import (
	"github.com/caicloud/nirvana"
	"github.com/caicloud/nirvana/definition"
	//	"github.com/caicloud/nirvana/log"
	//	"tfjob-admin/pkg/api/middlewares/logger"
)

var descriptors = []definition.Descriptor{}

func register(ds ...definition.Descriptor) {
	descriptors = append(descriptors, ds...)
}

func Initialize(s *nirvana.Config) {
	//	middlewars := []definition.Middleware{
	//		logger.New(log.DefaultLogger()),
	//	}

	s.Configure(nirvana.Descriptor(descriptors...))
	/*	v1 := definition.DescriptorFor("/api/v1alpha2", "v1alpha2 API").
			Middleware(middleware...).
			Consume(definition.MIMEJSON).
			Produce(definition.MIMEJSON).
			Descriptor(descriptors...)

		s.Configure(nirvana.Descriptor(v1))
	*/
}
