package rest

import (
	"net/http"

	"github.com/skyrocketOoO/AuthNet/domain"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	usecase domain.Usecase
}

func NewDelivery(usecase domain.Usecase) *Delivery {
	return &Delivery{
		usecase: usecase,
	}
}

// @Summary Check the server started
// @Accept json
// @Produce json
// @Success 200 {obj} domain.Response
// @Router /ping [get]
func (d *Delivery) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, domain.Response{Msg: "pong"})
}

// @Summary Check the server healthy
// @Accept json
// @Produce json
// @Success 200 {obj} domain.Response
// @Failure 503 {obj} domain.Response
// @Router /healthy [get]
func (d *Delivery) Healthy(c *gin.Context) {
	// do something check
	if err := d.usecase.Healthy(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, domain.Response{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response{Msg: "healthy"})
}

// @Summary Query edges based on parameters
// @Description Query edges based on specified parameters.
// @Tags Edge
// @Accept json
// @Produce json
// @Param obj-namespace query string false "Obj Namespace"
// @Param obj-name query string false "Obj Name"
// @Param edge query string false "Edge"
// @Param sbj-namespace query string false "Sbj Namespace"
// @Param sbj-name query string false "Sbj Name"
// @Param sbj-edge query string false "Sbj Edge"
// @Param page-token query string false "Page token"
// @Param page-size query string false "Page size"
// @Success 200 {obj} delivery.Get.respBody
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/ [get]
func (h *Delivery) Get(c *gin.Context) {
	edge := domain.Edge{
		ObjNs:   c.Query("obj_ns"),
		ObjName: c.Query("obj_name"),
		ObjRel:  c.Query("obj_rel"),
		SbjNs:   c.Query("sbj_ns"),
		SbjName: c.Query("sbj_name"),
		SbjRel:  c.Query("sbj_rel"),
	}
	queryMode := c.Query("query_mode") == "true"
	edges, err := h.usecase.Get(c.Request.Context(), edge, queryMode)
	if err != nil {
		if err.Error() == (domain.ErrRecordNotFound{}).Error() {
			c.JSON(http.StatusNotFound, domain.Response{Msg: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, domain.Response{
				Msg: err.Error(),
			})
		}
		return
	}
	type respBody struct {
		Edges []domain.Edge `json:"edges"`
	}
	c.JSON(http.StatusOK, respBody{
		Edges: edges,
	})
}

// @Summary Create a new edge
// @Description Create a new edge based on the provided JSON payload.
// @Tags Edge
// @Accept json
// @Produce json
// @Param edge body delivery.Create.requestBody true "Edge obj to be created"
// @Success 200
// @Failure 400 {obj} domain.ErrResponse
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/ [post]
func (h *Delivery) Create(c *gin.Context) {
	type requestBody struct {
		Edge domain.Edge `json:"edge"`
	}
	reqBody := requestBody{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	if err := h.usecase.Create(c.Request.Context(), reqBody.Edge); err != nil {
		if _, ok := err.(domain.ErrGraphCycle); ok {
			c.JSON(http.StatusBadRequest, domain.Response{
				Msg: err.Error(),
			})
		} else if _, ok := err.(domain.ErrBodyAttribute); ok {
			c.JSON(http.StatusBadRequest, domain.Response{
				Msg: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, domain.Response{
				Msg: err.Error(),
			})
		}
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Delete a edge
// @Description Delete a edge based on the provided JSON payload.
// @Tags Edge
// @Accept json
// @Produce json
// @Param edge body delivery.Delete.requestBody true "Edge obj to be deleted"
// @Success 200
// @Failure 400 {obj} domain.ErrResponse
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/ [delete]
func (h *Delivery) Delete(c *gin.Context) {
	type requestBody struct {
		Edge      domain.Edge `json:"edge"`
		QueryMode bool        `json:"query_mode"`
	}
	reqBody := requestBody{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Msg: err.Error(),
		})
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), reqBody.Edge,
		reqBody.QueryMode); err != nil {
		if err.Error() == (domain.ErrRecordNotFound{}).Error() {
			c.JSON(http.StatusNotFound, domain.Response{Msg: err.Error()})
			return
		}
		if _, ok := err.(domain.ErrBodyAttribute); ok {
			c.JSON(http.StatusBadRequest, domain.Response{
				Msg: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Clear all edges
// @Description Clear all edges in the system
// @Tags Edge
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/clear-all-edges [post]
func (h *Delivery) ClearAll(c *gin.Context) {
	err := h.usecase.ClearAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Check if a edge link exists
// @Description Check if a edge link exists between two entities
// @Tags Edge
// @Accept json
// @Produce json
// @Param edge body delivery.Check.requestBody true "comment"
// @Success 200
// @Failure 400 {obj} domain.ErrResponse
// @Failure 403
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/check [post]
func (h *Delivery) CheckAuth(c *gin.Context) {
	type requestBody struct {
		Sbj        domain.Vertex     `json:"sbj" binding:"required"`
		Obj        domain.Vertex     `json:"obj" binding:"required"`
		SearchCond domain.SearchCond `json:"search_cond"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	ok, err := h.usecase.CheckAuth(c.Request.Context(), body.Sbj, body.Obj,
		body.SearchCond)
	if err != nil {
		if _, ok := err.(domain.ErrBodyAttribute); ok {
			c.JSON(http.StatusBadRequest, domain.Response{
				Msg: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	if !ok {
		c.Status(http.StatusForbidden)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Get all edges for a given obj
// @Description Get all edges for a given obj specified by namespace, name, and edge
// @Tags Edge
// @Accept json
// @Produce json
// @Param sbj body delivery.GetAllObjEdges.requestBody true "Obj information (namespace, name, edge)"
// @Success 200 {obj} domain.DataResponse "All edges for the specified obj"
// @Failure 400 {obj} domain.ErrResponse
// @Failure 403
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/get-all-obj-edges [post]
func (h *Delivery) GetObjAuths(c *gin.Context) {
	type requestBody struct {
		Sbj         domain.Vertex      `json:"sbj" binding:"required"`
		SearchCond  domain.SearchCond  `json:"search_cond"`
		CollectCond domain.CollectCond `json:"collect_cond"`
		MaxDepth    int                `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	vertices, err := h.usecase.GetObjAuths(
		c.Request.Context(),
		domain.Vertex(body.Sbj),
		body.SearchCond,
		body.CollectCond,
		body.MaxDepth,
	)
	if err != nil {
		if _, ok := err.(domain.ErrBodyAttribute); ok {
			c.JSON(http.StatusBadRequest, domain.Response{
				Msg: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	type Response struct {
		Vertices []domain.Vertex `json:"vertices"`
	}
	c.JSON(http.StatusOK, Response{
		Vertices: vertices,
	})
}

// @Summary Get all edges for a given sbj
// @Description Get all edges for a given sbj specified by namespace, name, and edge
// @Tags Edge
// @Accept json
// @Produce json
// @Param obj body delivery.GetAllSbjEdges.requestBody true "Sbj information (namespace, name, edge)"
// @Success 200 {obj} domain.DataResponse "All edges for the specified sbj"
// @Failure 400 {obj} domain.ErrResponse
// @Failure 403
// @Failure 500 {obj} domain.ErrResponse
// @Router /edge/get-all-sbj-edges [post]
func (h *Delivery) GetSbjsWhoHasAuth(c *gin.Context) {
	type requestBody struct {
		Obj         domain.Vertex      `json:"obj" binding:"required"`
		SearchCond  domain.SearchCond  `json:"search_cond"`
		CollectCond domain.CollectCond `json:"collect_cond"`
		MaxDepth    int                `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Msg: err.Error(),
		})
		return
	}

	vertices, err := h.usecase.GetSbjsWhoHasAuth(
		c.Request.Context(),
		domain.Vertex(body.Obj),
		body.SearchCond,
		body.CollectCond,
		body.MaxDepth,
	)
	if err != nil {
		if _, ok := err.(domain.ErrBodyAttribute); ok {
			c.JSON(http.StatusBadRequest, domain.Response{
				Msg: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	type Response struct {
		Vertices []domain.Vertex `json:"vertices"`
	}
	c.JSON(http.StatusOK, Response{
		Vertices: vertices,
	})
}

func (h *Delivery) GetTree(c *gin.Context) {
	type requestBody struct {
		Sbj      domain.Vertex `json:"sbj" binding:"required"`
		MaxDepth int           `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	tree, err := h.usecase.GetTree(
		c.Request.Context(),
		body.Sbj,
		body.MaxDepth,
	)
	if err != nil {
		if _, ok := err.(domain.ErrBodyAttribute); ok {
			c.JSON(http.StatusBadRequest, domain.Response{
				Msg: err.Error(),
			})
			return
		} else if _, ok := err.(domain.ErrRecordNotFound); ok {
			c.JSON(http.StatusNotFound, domain.Response{
				Msg: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	type response struct {
		Tree domain.TreeNode `json:"tree"`
	}
	c.JSON(http.StatusOK, response{
		Tree: *tree,
	})
}

func (h *Delivery) SeeTree(c *gin.Context) {
	type requestBody struct {
		Sbj      domain.Vertex `json:"sbj" binding:"required"`
		MaxDepth int           `json:"max_depth"`
	}
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Msg: err.Error(),
		})
		return
	}
	graph, err := h.usecase.SeeTree(
		c.Request.Context(),
		body.Sbj,
		body.MaxDepth,
	)
	if err != nil {
		if _, ok := err.(domain.ErrBodyAttribute); ok {
			c.JSON(http.StatusBadRequest, domain.Response{
				Msg: err.Error(),
			})
			return
		} else if _, ok := err.(domain.ErrRecordNotFound); ok {
			c.JSON(http.StatusNotFound, domain.Response{
				Msg: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{
			Msg: err.Error(),
		})
		return
	}

	graph.Render(c.Writer)
}
