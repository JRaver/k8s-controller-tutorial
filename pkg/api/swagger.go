package api

// @Summary List all frontend pages
// @Description Get a list of all frontend pages
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDocList
// @Router /api/frontendpages [get]
// @Security ApiKeyAuth
// @Param namespace query string false "Namespace to filter by"
func SwaggerListFrontendPages() {}

// @Summary Get a frontend page
// @Description Get a frontend page by name
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDoc
// @Router /api/frontendpages/{name} [get]
// @Security ApiKeyAuth
// @Param name path string true "Name of the frontend page"
func SwaggerGetFrontendPage() {}

// @Summary Create a frontend page
// @Description Create a new frontend page
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDoc
// @Router /api/frontendpages [post]
// @Security ApiKeyAuth
// @Param frontendpage body FrontendPageDoc true "Frontend page to create"
func SwaggerCreateFrontendPage() {}

// @Summary Update a frontend page
// @Description Update a frontend page by name
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDoc
// @Router /api/frontendpages/{name} [put]
// @Security ApiKeyAuth
// @Param name path string true "Name of the frontend page"
func SwaggerUpdateFrontendPage() {}

// @Summary Delete a frontend page
// @Description Delete a frontend page by name
// @Tags frontendpages
// @Accept json
// @Produce json
// @Success 200 {object} FrontendPageDoc
// @Router /api/frontendpages/{name} [delete]
// @Security ApiKeyAuth
// @Param name path string true "Name of the frontend page"
func SwaggerDeleteFrontendPage() {}
