package validator

const (
	TagVBase                = "validator"
	TagSeparator            = "="
	MultiValidatorSeparator = ";"
	MultiChoiceSeparator    = ","
	VRequired               = "required"      // Example MyField string `validator:"required"`
	VGreaterThan            = "gt="           // Example MyField int `validator:"gt=1"`
	VGreaterThanOrEqual     = "gte="          // Example MyField int `validator:"gte=2"`
	VExactStr               = "exact="        // Example MyField string `validator:"exact=specificString"`
	VEqualTo                = "eq="           // Example MyField int `validator:"eq=3.2"`
	VLessThanOrEqual        = "lte="          // Example MyField int `validator:"lte=5"`
	VLessThan               = "lt="           // Example MyField int `validator:"lt=100"`
	VOneOf                  = "oneof="        // Example MyChoice string `validator:"oneof=cat,dog,raccoon,possum"`
	VEmail                  = "email"         // Example MyEmail string `validator:"email"`
	vInvalid                = "bad tag"       // Intentionally not exported. Used for switch/case control
	vNoTag                  = "no tag"        // Intentionally not exported. Used for switch/case control
	vMultipleTags           = "multiple tags" // Intentionally not exported. Used for switch/case control
)
