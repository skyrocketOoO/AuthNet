package mongo2

type Vertex struct {
	Ns   string `bson:"ns"`
	Name string `bson:"name"`
	Rel  string `bson:"rel"`
}

type Edge struct {
	UID string `bson:"u_id"`
	VID string `bson:"v_id"`
}
