package descriptor

import (
	//	"fmt"
	"github.com/caicloud/nirvana/definition"
	"tfjob-admin/pkg/api/v1/handlers"
)

func init() {
	register(app)
}

var app = definition.Descriptor{
	Path:        "/api/v1alpha2",
	Description: "Application API",
	Children: []definition.Descriptor{
		{
			Path:        "/clusters/{cid}/partitions/{partition}/tfjobs",
			Definitions: []definition.Definition{listjobs, createjobs},
		},
		{
			Path:        "/clusters/{cid}/paritions/{partition}/tfjobs/{jobid}",
			Definitions: []definition.Definition{getjobs, updatejobs, deletejobs},
		},
	},
}
var listjobs = definition.Definition{
	Method:      definition.List,
	Description: "list tfjobs",
	Function:    handlers.ListJobs,
	Parameters: []definition.Parameter{
		{
			Source:      definition.Path,
			Name:        "cid",
			Description: "cid name",
		},
		{
			Source:      definition.Path,
			Name:        "paritions",
			Description: "paritions name",
		},
	},
	Results: definition.DataErrorResults("all app list"),
}
var createjobs = definition.Definition{
	Method:      definition.Create,
	Description: "create a tfjob",
	Function:    handlers.CreateJob,
	Parameters: []definition.Parameter{
		{
			Source:      definition.Path,
			Name:        "cid",
			Description: "cid name",
		},
		{
			Source:      definition.Path,
			Name:        "paritions",
			Description: "parition name",
		},
		{
			Source:      definition.Body,
			Name:        "valuesReader",
			Description: "tfjob config",
		},
	},
	Results: definition.DataErrorResults("get a tfjob"),
}
var deletejobs = definition.Definition{
	Method:      definition.Delete,
	Description: "delete a tfjobs",
	Function:    handlers.DeleteJob,
	Parameters: []definition.Parameter{
		{
			Source:      definition.Path,
			Name:        "cid",
			Description: "cid name",
		},
		{
			Source:      definition.Path,
			Name:        "paritions",
			Description: "paritions name",
		},
		{
			Source:      definition.Path,
			Name:        "jobid",
			Description: "tfjob id",
		},
	},
	Results: definition.DataErrorResults("delete a tfjob"),
}
var updatejobs = definition.Definition{
	Method:      definition.Update,
	Description: "updata a tfjobs",
	Function:    handlers.UpdateJob,
	Parameters: []definition.Parameter{
		{
			Source:      definition.Path,
			Name:        "cid",
			Description: "cid",
		},
		{
			Source:      definition.Path,
			Name:        "paritions",
			Description: "paritions name",
		},
		{
			Source:      definition.Path,
			Name:        "jobid",
			Description: "tfjob id",
		},
		{
			Source:      definition.Body,
			Name:        "valueReader",
			Description: "tfjob config",
		},
	},
	Results: definition.DataErrorResults("update a tfjob"),
}
var getjobs = definition.Definition{
	Method:      definition.Get,
	Description: "get a tfjob",
	Function:    handlers.GetJob,
	Parameters: []definition.Parameter{
		{
			Source:      definition.Path,
			Name:        "cid",
			Description: "cid",
		},
		{
			Source:      definition.Path,
			Name:        "paritions",
			Description: "paritions name",
		},
		{
			Source:      definition.Path,
			Name:        "jobid",
			Description: "tfjob id",
		},
	},
	Results: definition.DataErrorResults("all app list"),
}
