# Test Coverage Audit

## Scope and Method

- Static inspection only; no code/tests/scripts/containers executed.
- Endpoint inventory source: `repo/cmd/server/main.go`.
- API test evidence sources: `repo/API_tests/*_test.go`, `repo/internal/integration/*_test.go`.
- Unit/frontend evidence sources: `repo/unit_tests/*_test.go`, `repo/internal/**/_test.go`, `repo/web/tests/*.test.js`, `repo/web/package.json`, `repo/web/vitest.config.js`.

## Project Type Detection

- README explicitly declares project type at top: `repo/README.md:3` -> **fullstack**.

## Backend Endpoint Inventory

Extracted endpoints (METHOD + PATH):

`GET /health`, `GET /metrics`, `GET /`, `GET /login`, `POST /login`, `GET /share/:token`, `GET /account/change-password`, `POST /account/change-password`, `POST /logout`, `GET /dashboard`, `GET /history`, `GET /materials`, `GET /materials/search`, `GET /materials/:id`, `POST /materials/:id/rate`, `POST /materials/:id/comments`, `POST /comments/:id/report`, `GET /favorites`, `POST /favorites`, `GET /favorites/:id`, `GET /favorites/:id/items`, `POST /favorites/:id/items`, `DELETE /favorites/:id/items/:materialID`, `GET /favorites/:id/share`, `GET /orders`, `GET /orders/cart`, `GET /orders/:id`, `POST /orders`, `POST /orders/:id/pay`, `POST /orders/:id/cancel`, `POST /orders/:id/returns`, `GET /returns`, `GET /inbox`, `GET /inbox/items`, `POST /inbox/:id/read`, `POST /inbox/read-all`, `GET /inbox/settings`, `POST /inbox/settings/dnd`, `POST /inbox/subscribe`, `POST /inbox/unsubscribe`, `GET /inbox/badge`, `GET /api/inbox/unread-count`, `GET /inbox/sse`, `GET /api/stats/:stat`, `GET /distribution`, `POST /distribution/issue`, `POST /distribution/return`, `POST /distribution/exchange`, `GET /distribution/reissue`, `POST /distribution/reissue`, `GET /distribution/ledger`, `GET /distribution/ledger/search`, `GET /distribution/custody/:scanID`, `GET /admin/orders`, `POST /admin/orders/:id/ship`, `POST /admin/orders/:id/deliver`, `GET /moderation`, `GET /moderation/items`, `POST /moderation/:id/approve`, `POST /moderation/:id/remove`, `GET /dashboard/instructor`, `GET /admin/returns`, `POST /admin/returns/:id/approve`, `POST /admin/returns/:id/reject`, `POST /admin/orders/:id/cancel`, `GET /courses`, `GET /courses/new`, `POST /courses`, `GET /courses/:id`, `POST /courses/:id/plan`, `POST /courses/:id/plan/:planID/approve`, `POST /courses/:id/sections`, `GET /dashboard/admin`, `GET /admin/materials/new`, `POST /admin/materials`, `GET /admin/materials/:id/edit`, `PUT /admin/materials/:id`, `DELETE /admin/materials/:id`, `GET /admin/users`, `GET /admin/users/new`, `POST /admin/users`, `GET /admin/users/:id`, `POST /admin/users/:id/role`, `POST /admin/users/:id/unlock`, `GET /admin/fields/:entity_type/:entity_id`, `POST /admin/fields/:entity_type/:entity_id`, `DELETE /admin/fields/:entity_type/:entity_id/:name`, `GET /admin/users/:id/fields`, `POST /admin/users/:id/fields`, `DELETE /admin/users/:id/fields/:name`, `GET /admin/duplicates`, `POST /admin/duplicates/merge`, `GET /admin/audit`, `GET /admin/audit/:entityType/:entityID`, `GET /analytics/map`, `GET /analytics/map/data`, `POST /analytics/map/compute`, `GET /analytics/map/buffer`, `GET /analytics/map/poi-density`, `GET /analytics/map/trajectory/:materialID`, `GET /analytics/map/regions`, `POST /analytics/map/regions/compute`, `GET /analytics/export/orders`, `GET /analytics/export/distribution`, `GET /analytics/kpi/:name`.

Total endpoints: **105**.

## API Test Mapping Table

| Endpoint | covered | test type | test files | evidence |
|---|---|---|---|---|
| `GET /health` | yes | true no-mock HTTP | `API_tests/auth_test.go` | `API_tests/auth_test.go:18` `TestAuth_Health` |
| `GET /metrics` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:295` `TestPermission_Metrics_AdminAllowed` |
| `GET /` | yes | true no-mock HTTP | `API_tests/missing_routes_test.go` | `API_tests/missing_routes_test.go:42` `TestRoot_Authenticated_Redirects` |
| `GET /login` | yes | true no-mock HTTP | `API_tests/missing_routes_test.go` | `API_tests/missing_routes_test.go:60` `TestLogin_Page_RendersForm` |
| `POST /login` | yes | true no-mock HTTP | `API_tests/auth_test.go` | `API_tests/auth_test.go:34` `TestAuth_Login_ValidCredentials` |
| `GET /share/:token` | yes | true no-mock HTTP | `API_tests/edge_cases_test.go` | `API_tests/edge_cases_test.go:41` `TestEdge_ShareLink_ValidToken` |
| `GET /account/change-password` | yes | true no-mock HTTP | `API_tests/account_test.go` | `API_tests/account_test.go:16` `TestAccount_ChangePasswordPage_ReturnsOK` |
| `POST /account/change-password` | yes | true no-mock HTTP | `API_tests/account_test.go` | `API_tests/account_test.go:40` `TestAccount_ChangePassword_Valid` |
| `POST /logout` | yes | true no-mock HTTP | `API_tests/auth_test.go` | `API_tests/auth_test.go:60` `TestAuth_Logout_ClearsSession` |
| `GET /dashboard` | yes | true no-mock HTTP | `API_tests/auth_test.go` | `API_tests/auth_test.go:72` `TestAuth_Dashboard_AuthenticatedStudent` |
| `GET /history` | yes | true no-mock HTTP | `API_tests/history_test.go` | `API_tests/history_test.go:17` `TestHistory_ReturnsOK` |
| `GET /materials` | yes | true no-mock HTTP | `API_tests/materials_test.go` | `API_tests/materials_test.go:20` `TestMaterials_List_ReturnsOK` |
| `GET /materials/search` | yes | true no-mock HTTP | `API_tests/materials_test.go` | `API_tests/materials_test.go:48` `TestMaterials_Search_WithQuery` |
| `GET /materials/:id` | yes | true no-mock HTTP | `API_tests/materials_test.go` | `API_tests/materials_test.go:34` `TestMaterials_Detail_ReturnsOKOrTemplate` |
| `POST /materials/:id/rate` | yes | true no-mock HTTP | `API_tests/materials_test.go` | `API_tests/materials_test.go:81` `TestMaterials_Rate_ValidStars` |
| `POST /materials/:id/comments` | yes | true no-mock HTTP | `API_tests/materials_test.go` | `API_tests/materials_test.go:97` `TestMaterials_AddComment_Valid` |
| `POST /comments/:id/report` | yes | true no-mock HTTP | `API_tests/edge_cases_test.go` | `API_tests/edge_cases_test.go:101` `TestEdge_ReportComment_CollapseAt3` |
| `GET /favorites` | yes | true no-mock HTTP | `API_tests/missing_routes_test.go` | `API_tests/missing_routes_test.go:95` `TestFavorites_Index_AuthenticatedReturnsOK` |
| `POST /favorites` | yes | true no-mock HTTP | `API_tests/favorites_extended_test.go` | `API_tests/favorites_extended_test.go:31` `TestFavorites_Detail_Owner` |
| `GET /favorites/:id` | yes | true no-mock HTTP | `API_tests/favorites_extended_test.go` | `API_tests/favorites_extended_test.go:48` `TestFavorites_Detail_Owner` |
| `GET /favorites/:id/items` | yes | true no-mock HTTP | `API_tests/favorites_extended_test.go` | `API_tests/favorites_extended_test.go:99` `TestFavorites_Items_ReturnsOK` |
| `POST /favorites/:id/items` | yes | true no-mock HTTP | `API_tests/materials_test.go` | `API_tests/materials_test.go:138` `TestMaterials_Favorites_AddItem` |
| `DELETE /favorites/:id/items/:materialID` | yes | true no-mock HTTP | `API_tests/missing_routes_test.go` | `API_tests/missing_routes_test.go:153` `TestFavorites_RemoveItem_OwnerSucceeds` |
| `GET /favorites/:id/share` | yes | true no-mock HTTP | `API_tests/favorites_extended_test.go` | `API_tests/favorites_extended_test.go:151` `TestFavorites_Share_ReturnsOK` |
| `GET /orders` | yes | true no-mock HTTP | `API_tests/orders_test.go` | `API_tests/orders_test.go:22` `TestOrders_List_ReturnsOK` |
| `GET /orders/cart` | yes | true no-mock HTTP | `API_tests/orders_test.go` | `API_tests/orders_test.go:106` `TestOrders_Cart_ReturnsOK` |
| `GET /orders/:id` | yes | true no-mock HTTP | `internal/integration/orders_test.go` | `internal/integration/orders_test.go:77` `TestOrderDetail` |
| `POST /orders` | yes | true no-mock HTTP | `API_tests/orders_test.go` | `API_tests/orders_test.go:36` `TestOrders_PlaceOrder_ValidItem` |
| `POST /orders/:id/pay` | yes | true no-mock HTTP | `API_tests/analytics_extended_test.go` | `API_tests/analytics_extended_test.go:219` `TestOrders_ConfirmPayment_Valid` |
| `POST /orders/:id/cancel` | yes | true no-mock HTTP | `API_tests/orders_test.go` | `API_tests/orders_test.go:92` `TestOrders_Cancel_PendingPayment` |
| `POST /orders/:id/returns` | yes | true no-mock HTTP | `API_tests/returns_test.go` | `API_tests/returns_test.go:82` `TestReturns_Submit_Valid` |
| `GET /returns` | yes | true no-mock HTTP | `API_tests/returns_test.go` | `API_tests/returns_test.go:29` `TestReturns_List_ReturnsOK` |
| `GET /inbox` | yes | true no-mock HTTP | `API_tests/edge_cases_test.go` | `API_tests/edge_cases_test.go:188` `TestEdge_Inbox_ReturnsOK` |
| `GET /inbox/items` | yes | true no-mock HTTP | `API_tests/inbox_extended_test.go` | `API_tests/inbox_extended_test.go:29` `TestInbox_Items_ReturnsOK` |
| `POST /inbox/:id/read` | yes | true no-mock HTTP | `API_tests/inbox_extended_test.go` | `API_tests/inbox_extended_test.go:57` `TestInbox_MarkRead_NonExistentID` |
| `POST /inbox/read-all` | yes | true no-mock HTTP | `API_tests/inbox_extended_test.go` | `API_tests/inbox_extended_test.go:86` `TestInbox_MarkAllRead_ReturnsOK` |
| `GET /inbox/settings` | yes | true no-mock HTTP | `API_tests/inbox_extended_test.go` | `API_tests/inbox_extended_test.go:241` `TestInbox_Settings_ReturnsOK` |
| `POST /inbox/settings/dnd` | yes | true no-mock HTTP | `API_tests/inbox_extended_test.go` | `API_tests/inbox_extended_test.go:131` `TestInbox_UpdateDND_ValidHours` |
| `POST /inbox/subscribe` | yes | true no-mock HTTP | `API_tests/inbox_extended_test.go` | `API_tests/inbox_extended_test.go:189` `TestInbox_Subscribe_ReturnsOK` |
| `POST /inbox/unsubscribe` | yes | true no-mock HTTP | `API_tests/inbox_extended_test.go` | `API_tests/inbox_extended_test.go:207` `TestInbox_Unsubscribe_ReturnsOK` |
| `GET /inbox/badge` | yes | true no-mock HTTP | `API_tests/edge_cases_test.go` | `API_tests/edge_cases_test.go:200` `TestEdge_Badge_ReturnsOK` |
| `GET /api/inbox/unread-count` | yes | true no-mock HTTP | `API_tests/inbox_extended_test.go` | `API_tests/inbox_extended_test.go:103` `TestInbox_UnreadCount_ReturnsOK` |
| `GET /inbox/sse` | yes | true no-mock HTTP | `internal/integration/distribution_test.go` | `internal/integration/distribution_test.go:387` `TestInboxSSE_RouteRegistered` |
| `GET /api/stats/:stat` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:442` `TestPermission_Stats_AdminCanAccessAllStats` |
| `GET /distribution` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:78` `TestPermission_Distribution_ClerkAllowed` |
| `POST /distribution/issue` | yes | true no-mock HTTP | `internal/integration/distribution_test.go` | `internal/integration/distribution_test.go:59` `TestIssueItems_Success` |
| `POST /distribution/return` | yes | true no-mock HTTP | `API_tests/distribution_extended_test.go` | `API_tests/distribution_extended_test.go:108` `TestDistribution_RecordReturn_WithApprovedRequest` |
| `POST /distribution/exchange` | yes | true no-mock HTTP | `API_tests/distribution_extended_test.go` | `API_tests/distribution_extended_test.go:126` `TestDistribution_RecordExchange_MissingFields` |
| `GET /distribution/reissue` | yes | true no-mock HTTP | `API_tests/distribution_extended_test.go` | `API_tests/distribution_extended_test.go:157` `TestDistribution_ReissueForm_ClerkAllowed` |
| `POST /distribution/reissue` | yes | true no-mock HTTP | `API_tests/distribution_extended_test.go` | `API_tests/distribution_extended_test.go:235` `TestDistribution_Reissue_ValidFlow` |
| `GET /distribution/ledger` | yes | true no-mock HTTP | `internal/integration/distribution_test.go` | `internal/integration/distribution_test.go:101` `TestLedger_ReturnsEntries` |
| `GET /distribution/ledger/search` | yes | true no-mock HTTP | `API_tests/distribution_extended_test.go` | `API_tests/distribution_extended_test.go:253` `TestDistribution_LedgerSearch_ClerkAllowed` |
| `GET /distribution/custody/:scanID` | yes | true no-mock HTTP | `API_tests/distribution_extended_test.go` | `API_tests/distribution_extended_test.go:325` `TestDistribution_CustodyChain_ExistingScan` |
| `GET /admin/orders` | yes | true no-mock HTTP | `API_tests/authorized_routes_test.go` | `API_tests/authorized_routes_test.go:66` `TestAdminOrders_WithOrders_ClerkSeesAll` |
| `POST /admin/orders/:id/ship` | yes | true no-mock HTTP | `API_tests/orders_test.go` | `API_tests/orders_test.go:250` `TestOrders_MarkShipped_ClerkAllowed` |
| `POST /admin/orders/:id/deliver` | yes | true no-mock HTTP | `API_tests/returns_test.go` | `API_tests/returns_test.go:269` `TestAdminOrders_Deliver_ClerkAllowed` |
| `GET /moderation` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:38` `TestPermission_Moderation_ModeratorAllowed` |
| `GET /moderation/items` | yes | true no-mock HTTP | `API_tests/missing_routes_test.go` | `API_tests/missing_routes_test.go:200` `TestModeration_Items_ModeratorAllowed` |
| `POST /moderation/:id/approve` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:265` `TestPermission_ApproveComment_ModeratorAllowed` |
| `POST /moderation/:id/remove` | yes | true no-mock HTTP | `internal/integration/moderation_test.go` | `internal/integration/moderation_test.go:105` `TestRemoveComment` |
| `GET /dashboard/instructor` | yes | true no-mock HTTP | `API_tests/missing_routes_test.go` | `API_tests/missing_routes_test.go:239` `TestDashboard_Instructor_InstructorAllowed` |
| `GET /admin/returns` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:159` `TestPermission_AdminReturns_InstructorAllowed` |
| `POST /admin/returns/:id/approve` | yes | true no-mock HTTP | `API_tests/returns_test.go` | `API_tests/returns_test.go:147` `TestReturns_Approve_InstructorCanApprove` |
| `POST /admin/returns/:id/reject` | yes | true no-mock HTTP | `API_tests/returns_test.go` | `API_tests/returns_test.go:185` `TestReturns_Reject_ManagerCanReject` |
| `POST /admin/orders/:id/cancel` | yes | true no-mock HTTP | `API_tests/returns_test.go` | `API_tests/returns_test.go:227` `TestAdminOrders_Cancel_InstructorAllowed` |
| `GET /courses` | yes | true no-mock HTTP | `API_tests/courses_test.go` | `API_tests/courses_test.go:29` `TestCourses_List_InstructorAllowed` |
| `GET /courses/new` | yes | true no-mock HTTP | `API_tests/courses_test.go` | `API_tests/courses_test.go:70` `TestCourses_NewForm_InstructorAllowed` |
| `POST /courses` | yes | true no-mock HTTP | `API_tests/courses_test.go` | `API_tests/courses_test.go:99` `TestCourses_Create_InstructorAllowed` |
| `GET /courses/:id` | yes | true no-mock HTTP | `API_tests/courses_test.go` | `API_tests/courses_test.go:173` `TestCourses_Detail_NotFound` |
| `POST /courses/:id/plan` | yes | true no-mock HTTP | `API_tests/courses_test.go` | `API_tests/courses_test.go:207` `TestCourses_AddPlanItem_InstructorAllowed` |
| `POST /courses/:id/plan/:planID/approve` | yes | true no-mock HTTP | `API_tests/courses_test.go` | `API_tests/courses_test.go:332` `TestCourses_ApprovePlanItem_AdminAllowed` |
| `POST /courses/:id/sections` | yes | true no-mock HTTP | `API_tests/courses_test.go` | `API_tests/courses_test.go:255` `TestCourses_AddSection_InstructorAllowed` |
| `GET /dashboard/admin` | yes | true no-mock HTTP | `internal/integration/auth_test.go` | `internal/integration/auth_test.go:255` `TestDefaultAdminLogin_SeedCredentialsWork` |
| `GET /admin/materials/new` | yes | true no-mock HTTP | `API_tests/admin_materials_test.go` | `API_tests/admin_materials_test.go:25` `TestAdminMaterials_NewForm_AdminAllowed` |
| `POST /admin/materials` | yes | true no-mock HTTP | `API_tests/materials_test.go` | `API_tests/materials_test.go:278` `TestAdmin_CreateMaterial_AdminAllowed` |
| `GET /admin/materials/:id/edit` | yes | true no-mock HTTP | `API_tests/admin_materials_test.go` | `API_tests/admin_materials_test.go:56` `TestAdminMaterials_EditForm_AdminAllowed` |
| `PUT /admin/materials/:id` | yes | true no-mock HTTP | `API_tests/admin_materials_test.go` | `API_tests/admin_materials_test.go:104` `TestAdminMaterials_Update_AdminAllowed` |
| `DELETE /admin/materials/:id` | yes | true no-mock HTTP | `API_tests/admin_materials_test.go` | `API_tests/admin_materials_test.go:164` `TestAdminMaterials_Delete_AdminAllowed` |
| `GET /admin/users` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:118` `TestPermission_AdminUsers_AdminAllowed` |
| `GET /admin/users/new` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:31` `TestAdminUsers_NewForm_AdminAllowed` |
| `POST /admin/users` | yes | true no-mock HTTP | `internal/integration/admin_test.go` | `internal/integration/admin_test.go:46` `TestCreateUser` |
| `GET /admin/users/:id` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:62` `TestAdminUsers_Profile_AdminAllowed` |
| `POST /admin/users/:id/role` | yes | true no-mock HTTP | `internal/integration/admin_test.go` | `internal/integration/admin_test.go:81` `TestUpdateRole` |
| `POST /admin/users/:id/unlock` | yes | true no-mock HTTP | `internal/integration/admin_test.go` | `internal/integration/admin_test.go:114` `TestUnlockUser` |
| `GET /admin/fields/:entity_type/:entity_id` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:117` `TestAdminFields_CustomFieldsPage_AdminAllowed` |
| `POST /admin/fields/:entity_type/:entity_id` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:151` `TestAdminFields_SetCustomField_AdminAllowed` |
| `DELETE /admin/fields/:entity_type/:entity_id/:name` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:209` `TestAdminFields_DeleteCustomField_AdminAllowed` |
| `GET /admin/users/:id/fields` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:242` `TestAdminFields_LegacyAlias_AdminAllowed` |
| `POST /admin/users/:id/fields` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:259` `TestAdminFields_LegacyPost_AdminAllowed` |
| `DELETE /admin/users/:id/fields/:name` | yes | true no-mock HTTP | `API_tests/authorized_routes_test.go` | `API_tests/authorized_routes_test.go:194` `TestAdminUserFields_LegacyDelete_AdminAllowed` |
| `GET /admin/duplicates` | yes | true no-mock HTTP | `internal/integration/admin_test.go` | `internal/integration/admin_test.go:144` `TestDuplicatesPage` |
| `POST /admin/duplicates/merge` | yes | true no-mock HTTP | `API_tests/missing_routes_test.go` | `API_tests/missing_routes_test.go:294` `TestDuplicates_Merge_AdminSucceeds` |
| `GET /admin/audit` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:308` `TestAdminAudit_Global_AdminAllowed` |
| `GET /admin/audit/:entityType/:entityID` | yes | true no-mock HTTP | `API_tests/admin_users_extended_test.go` | `API_tests/admin_users_extended_test.go:280` `TestAdminAudit_Entity_AdminAllowed` |
| `GET /analytics/map` | yes | true no-mock HTTP | `API_tests/analytics_extended_test.go` | `API_tests/analytics_extended_test.go:64` `TestAnalytics_MapPage_AdminAllowed` |
| `GET /analytics/map/data` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:335` `TestPermission_MapData_AdminAllowed` |
| `POST /analytics/map/compute` | yes | true no-mock HTTP | `API_tests/authorized_routes_test.go` | `API_tests/authorized_routes_test.go:82` `TestAnalytics_ComputeGrid_AdminAllowed` |
| `GET /analytics/map/buffer` | yes | true no-mock HTTP | `API_tests/authorized_routes_test.go` | `API_tests/authorized_routes_test.go:131` `TestAnalytics_BufferQuery_AdminAllowed` |
| `GET /analytics/map/poi-density` | yes | true no-mock HTTP | `API_tests/analytics_extended_test.go` | `API_tests/analytics_extended_test.go:93` `TestAnalytics_POIDensity_AdminAllowed` |
| `GET /analytics/map/trajectory/:materialID` | yes | true no-mock HTTP | `API_tests/analytics_extended_test.go` | `API_tests/analytics_extended_test.go:123` `TestAnalytics_Trajectory_AdminAllowed` |
| `GET /analytics/map/regions` | yes | true no-mock HTTP | `API_tests/analytics_extended_test.go` | `API_tests/analytics_extended_test.go:153` `TestAnalytics_Regions_AdminAllowed` |
| `POST /analytics/map/regions/compute` | yes | true no-mock HTTP | `API_tests/analytics_extended_test.go` | `API_tests/analytics_extended_test.go:170` `TestAnalytics_ComputeRegions_AdminAllowed` |
| `GET /analytics/export/orders` | yes | true no-mock HTTP | `API_tests/permissions_test.go` | `API_tests/permissions_test.go:142` `TestPermission_Analytics_AdminAllowed` |
| `GET /analytics/export/distribution` | yes | true no-mock HTTP | `API_tests/missing_routes_test.go` | `API_tests/missing_routes_test.go:354` `TestAnalytics_ExportDistribution_AdminAllowed` |
| `GET /analytics/kpi/:name` | yes | true no-mock HTTP | `API_tests/analytics_extended_test.go` | `API_tests/analytics_extended_test.go:32` `TestAnalytics_KPI_AdminAllowed` |

## API Test Classification

1. **True No-Mock HTTP**
- API tests build real Fiber app + repositories + services + handlers and send HTTP via `app.Test` (`repo/API_tests/helpers_test.go:40`, `repo/API_tests/helpers_test.go:382`, `repo/internal/integration/helpers_test.go:57`, `repo/internal/integration/helpers_test.go:459`).

2. **HTTP with Mocking**
- None detected in backend API/integration tests.

3. **Non-HTTP (unit/integration without HTTP)**
- `repo/unit_tests/*_test.go`.
- `repo/internal/services/*_test.go`.
- `repo/internal/repository/*_test.go`.
- `repo/internal/config/config_test.go`.
- `repo/internal/scheduler/scheduler_test.go`.
- `repo/internal/middleware/auth_rbac_test.go`, `repo/internal/middleware/ratelimit_test.go`.
- `repo/internal/crypto/*_test.go`, `repo/internal/db/migration005_test.go`, `repo/internal/auth/credentials_integrity_test.go`.

## Mock Detection Rules Check

- Backend API/integration test files contain no `jest.mock`, `vi.mock`, `sinon.stub`, gomock/testify/sqlmock indicators.
- Frontend unit tests use `vi.fn` mocks (`repo/web/tests/map.test.js`, `repo/web/tests/app.test.js`) for browser primitives; these are frontend unit tests, not API tests.

## Coverage Summary

- Total endpoints: **105**.
- Endpoints with HTTP tests: **105**.
- Endpoints with true no-mock tests: **105**.
- HTTP coverage %: **100%**.
- True API coverage %: **100%**.

## Unit Test Summary

### Backend Unit Tests

- Controllers: no direct handler-unit test files found (`repo/internal/handlers/*_test.go` absent).
- Services covered: auth, materials, orders, distribution, messaging, analytics, admin, courses, moderation (`repo/internal/services/*.go` test files present).
- Repositories covered: users, orders, materials, engagement, distribution, messaging, analytics, admin credentials (`repo/internal/repository/*_test.go`).
- Auth/guards/middleware covered: `repo/internal/middleware/auth_rbac_test.go`, `repo/internal/middleware/ratelimit_test.go`, plus unit auth tests in `repo/unit_tests/auth_test.go`.
- Config/scheduler now covered: `repo/internal/config/config_test.go`, `repo/internal/scheduler/scheduler_test.go`.

Important backend modules not directly unit-tested:
- `repo/internal/handlers/*.go` (no dedicated handler-unit tests).

### Frontend Unit Tests (STRICT REQUIREMENT)

- Frontend test files: `repo/web/tests/app.test.js`, `repo/web/tests/map.test.js`.
- Framework/tools detected: Vitest + jsdom (`repo/web/package.json:7-14`, `repo/web/vitest.config.js:4-7`).
- Frontend modules covered:
  - `repo/web/static/js/app.js` (toast/confirm helpers, DOM behavior, escaping)
  - `repo/web/static/js/map.js` (Leaflet init, fetch flow, DOM wiring, escaping)
- Important frontend modules not directly unit-tested:
  - No explicit unit tests for vendored libs under `repo/web/static/js/*.min.js` (acceptable for third-party, but still untested in-repo behavior around integration points).

**Frontend unit tests: PRESENT**

Strict failure rule (fullstack + missing frontend tests): **not triggered**.

### Cross-Layer Observation

- Balance improved: backend API/integration + backend unit + frontend unit are all present.
- Remaining imbalance: no dedicated handler-unit layer, but API/integration tests partially compensate.

## API Observability Check

- Strong: tests explicitly name method/path and include request body/query in many cases (`repo/API_tests/authorized_routes_test.go`, `repo/API_tests/missing_routes_test.go`, `repo/API_tests/orders_test.go`).
- Mixed: some tests focus mainly on status codes (permission matrix), with fewer deep response payload assertions.
- Verdict: **good observability with moderate assertion-depth variability**.

## Tests Check

- Success paths: broad coverage across auth, materials, favorites, orders, returns, moderation, admin, analytics, distribution.
- Failure/validation: broad coverage (permissions, invalid input, missing fields, CSRF, auth).
- Edge cases: present (`repo/API_tests/edge_cases_test.go`).
- Auth/permissions: extensive matrix (`repo/API_tests/permissions_test.go`).
- Integration boundaries: present with real app wiring (`repo/internal/integration/helpers_test.go`).

`run_tests.sh` static check:
- Docker wrapper for tests is present (`repo/run_tests.sh:23-50`) -> Docker-contained execution path.
- Script still performs `npm install` inside container (`repo/run_tests.sh:168`) -> runtime install step exists (note for strict policy).

## Test Coverage Score (0–100)

**89/100**

## Score Rationale

- + Endpoint coverage is complete (105/105), with no-mock HTTP wiring.
- + Backend unit layers expanded (config/scheduler/middleware/service/repo).
- + Frontend unit tests now present with direct module-level evidence.
- - Some API assertions remain shallow (status-heavy vs payload-depth).
- - No dedicated handler-unit test layer.

## Key Gaps

- Handler-unit test gap: `repo/internal/handlers/*_test.go` missing.
- Assertion-depth gap: several permission tests validate status only.

## Confidence & Assumptions

- Confidence: **high** for route inventory and static method/path mapping.
- Assumption: endpoint “covered” classification is based on visible static request evidence and test setup wiring; no runtime execution performed.

## Test Coverage Final Verdict

**PASS**


# README Audit

## Hard Gate Failures

- None.

## High Priority Issues

- None.

## Medium Priority Issues

- `run_tests.sh` still executes `npm install` at runtime inside Docker (`repo/run_tests.sh:168`); this is outside README but affects strict reproducibility posture.

## Low Priority Issues

- Admin credential is intentionally dynamic (`<SECURITY.temporary_password>`) and must be fetched from logs; documented clearly but operationally less convenient than fixed seeded credentials.

## Hard Gate Checks That Pass

- README exists at required location: `repo/README.md`.
- Formatting/readability: clean markdown structure and tables.
- Startup instruction includes required command: `docker-compose up` (`repo/README.md:29`).
- Access method is explicit: `http://localhost:3000` (`repo/README.md:34`).
- Verification is explicit with curl and UI flow (`repo/README.md:88-131`).
- Environment rules in README now avoid local install workflows (previous optional local section removed).
- Demo credentials include username + email + password + role table for all roles, including admin (`repo/README.md:74-85`).
- Project type declared at top (`repo/README.md:3`).

## README Verdict

**PASS**


# Combined Final Verdicts

- **Test Coverage Audit:** PASS
- **README Audit:** PASS
