/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(() => {
var exports = {};
exports.id = "pages/_app";
exports.ids = ["pages/_app"];
exports.modules = {

/***/ "./context/AuthContext.tsx":
/*!*********************************!*\
  !*** ./context/AuthContext.tsx ***!
  \*********************************/
/***/ ((module, __webpack_exports__, __webpack_require__) => {

"use strict";
eval("__webpack_require__.a(module, async (__webpack_handle_async_dependencies__, __webpack_async_result__) => { try {\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   AuthProvider: () => (/* binding */ AuthProvider),\n/* harmony export */   useAuth: () => (/* binding */ useAuth)\n/* harmony export */ });\n/* harmony import */ var react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react/jsx-dev-runtime */ \"react/jsx-dev-runtime\");\n/* harmony import */ var react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! react */ \"react\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(react__WEBPACK_IMPORTED_MODULE_1__);\n/* harmony import */ var axios__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! axios */ \"axios\");\n/* harmony import */ var js_cookie__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! js-cookie */ \"js-cookie\");\nvar __webpack_async_dependencies__ = __webpack_handle_async_dependencies__([axios__WEBPACK_IMPORTED_MODULE_2__, js_cookie__WEBPACK_IMPORTED_MODULE_3__]);\n([axios__WEBPACK_IMPORTED_MODULE_2__, js_cookie__WEBPACK_IMPORTED_MODULE_3__] = __webpack_async_dependencies__.then ? (await __webpack_async_dependencies__)() : __webpack_async_dependencies__);\n\n\n\n\nconst AuthContext = /*#__PURE__*/ (0,react__WEBPACK_IMPORTED_MODULE_1__.createContext)(undefined);\nconst AuthProvider = ({ children })=>{\n    const [user, setUser] = (0,react__WEBPACK_IMPORTED_MODULE_1__.useState)(null);\n    const [token, setToken] = (0,react__WEBPACK_IMPORTED_MODULE_1__.useState)(null);\n    const [loading, setLoading] = (0,react__WEBPACK_IMPORTED_MODULE_1__.useState)(true);\n    const API_URL = \"http://localhost:8080\" || 0;\n    (0,react__WEBPACK_IMPORTED_MODULE_1__.useEffect)(()=>{\n        const storedToken = js_cookie__WEBPACK_IMPORTED_MODULE_3__[\"default\"].get(\"auth_token\");\n        if (storedToken) {\n            setToken(storedToken);\n            fetchUser(storedToken);\n        } else {\n            setLoading(false);\n        }\n    }, []);\n    const fetchUser = async (authToken)=>{\n        try {\n            const response = await axios__WEBPACK_IMPORTED_MODULE_2__[\"default\"].get(`${API_URL}/api/users/me`, {\n                headers: {\n                    Authorization: `Bearer ${authToken}`\n                }\n            });\n            setUser(response.data);\n        } catch (error) {\n            console.error(\"Failed to fetch user:\", error);\n            js_cookie__WEBPACK_IMPORTED_MODULE_3__[\"default\"].remove(\"auth_token\");\n            setToken(null);\n        } finally{\n            setLoading(false);\n        }\n    };\n    const login = async (email, password)=>{\n        const response = await axios__WEBPACK_IMPORTED_MODULE_2__[\"default\"].post(`${API_URL}/api/auth/login`, {\n            email,\n            password\n        });\n        const { token: newToken, user: newUser } = response.data;\n        js_cookie__WEBPACK_IMPORTED_MODULE_3__[\"default\"].set(\"auth_token\", newToken, {\n            expires: 7\n        });\n        setToken(newToken);\n        setUser(newUser);\n    };\n    const register = async (email, password, name)=>{\n        const response = await axios__WEBPACK_IMPORTED_MODULE_2__[\"default\"].post(`${API_URL}/api/auth/register`, {\n            email,\n            password,\n            name\n        });\n        const { token: newToken, user: newUser } = response.data;\n        js_cookie__WEBPACK_IMPORTED_MODULE_3__[\"default\"].set(\"auth_token\", newToken, {\n            expires: 7\n        });\n        setToken(newToken);\n        setUser(newUser);\n    };\n    const logout = ()=>{\n        js_cookie__WEBPACK_IMPORTED_MODULE_3__[\"default\"].remove(\"auth_token\");\n        setToken(null);\n        setUser(null);\n    };\n    return /*#__PURE__*/ (0,react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxDEV)(AuthContext.Provider, {\n        value: {\n            user,\n            token,\n            loading,\n            login,\n            register,\n            logout,\n            isAuthenticated: !!user\n        },\n        children: children\n    }, void 0, false, {\n        fileName: \"/Users/bruno/infra/besend/frontend/context/AuthContext.tsx\",\n        lineNumber: 79,\n        columnNumber: 5\n    }, undefined);\n};\nconst useAuth = ()=>{\n    const context = (0,react__WEBPACK_IMPORTED_MODULE_1__.useContext)(AuthContext);\n    if (!context) {\n        throw new Error(\"useAuth must be used within AuthProvider\");\n    }\n    return context;\n};\n\n__webpack_async_result__();\n} catch(e) { __webpack_async_result__(e); } });//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9jb250ZXh0L0F1dGhDb250ZXh0LnRzeCIsIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7QUFBNkU7QUFDcEQ7QUFDTTtBQW1CL0IsTUFBTU8sNEJBQWNOLG9EQUFhQSxDQUE4Qk87QUFFeEQsTUFBTUMsZUFBd0QsQ0FBQyxFQUFFQyxRQUFRLEVBQUU7SUFDaEYsTUFBTSxDQUFDQyxNQUFNQyxRQUFRLEdBQUdULCtDQUFRQSxDQUFjO0lBQzlDLE1BQU0sQ0FBQ1UsT0FBT0MsU0FBUyxHQUFHWCwrQ0FBUUEsQ0FBZ0I7SUFDbEQsTUFBTSxDQUFDWSxTQUFTQyxXQUFXLEdBQUdiLCtDQUFRQSxDQUFDO0lBRXZDLE1BQU1jLFVBQVVDLHVCQUErQixJQUFJO0lBRW5EZCxnREFBU0EsQ0FBQztRQUNSLE1BQU1pQixjQUFjZixxREFBVyxDQUFDO1FBQ2hDLElBQUllLGFBQWE7WUFDZlAsU0FBU087WUFDVEUsVUFBVUY7UUFDWixPQUFPO1lBQ0xMLFdBQVc7UUFDYjtJQUNGLEdBQUcsRUFBRTtJQUVMLE1BQU1PLFlBQVksT0FBT0M7UUFDdkIsSUFBSTtZQUNGLE1BQU1DLFdBQVcsTUFBTXBCLGlEQUFTLENBQUMsQ0FBQyxFQUFFWSxRQUFRLGFBQWEsQ0FBQyxFQUFFO2dCQUMxRFMsU0FBUztvQkFBRUMsZUFBZSxDQUFDLE9BQU8sRUFBRUgsVUFBVSxDQUFDO2dCQUFDO1lBQ2xEO1lBQ0FaLFFBQVFhLFNBQVNHLElBQUk7UUFDdkIsRUFBRSxPQUFPQyxPQUFPO1lBQ2RDLFFBQVFELEtBQUssQ0FBQyx5QkFBeUJBO1lBQ3ZDdkIsd0RBQWMsQ0FBQztZQUNmUSxTQUFTO1FBQ1gsU0FBVTtZQUNSRSxXQUFXO1FBQ2I7SUFDRjtJQUVBLE1BQU1nQixRQUFRLE9BQU9DLE9BQWVDO1FBQ2xDLE1BQU1ULFdBQVcsTUFBTXBCLGtEQUFVLENBQUMsQ0FBQyxFQUFFWSxRQUFRLGVBQWUsQ0FBQyxFQUFFO1lBQUVnQjtZQUFPQztRQUFTO1FBQ2pGLE1BQU0sRUFBRXJCLE9BQU91QixRQUFRLEVBQUV6QixNQUFNMEIsT0FBTyxFQUFFLEdBQUdaLFNBQVNHLElBQUk7UUFDeER0QixxREFBVyxDQUFDLGNBQWM4QixVQUFVO1lBQUVHLFNBQVM7UUFBRTtRQUNqRHpCLFNBQVNzQjtRQUNUeEIsUUFBUXlCO0lBQ1Y7SUFFQSxNQUFNRyxXQUFXLE9BQU9QLE9BQWVDLFVBQWtCTztRQUN2RCxNQUFNaEIsV0FBVyxNQUFNcEIsa0RBQVUsQ0FBQyxDQUFDLEVBQUVZLFFBQVEsa0JBQWtCLENBQUMsRUFBRTtZQUFFZ0I7WUFBT0M7WUFBVU87UUFBSztRQUMxRixNQUFNLEVBQUU1QixPQUFPdUIsUUFBUSxFQUFFekIsTUFBTTBCLE9BQU8sRUFBRSxHQUFHWixTQUFTRyxJQUFJO1FBQ3hEdEIscURBQVcsQ0FBQyxjQUFjOEIsVUFBVTtZQUFFRyxTQUFTO1FBQUU7UUFDakR6QixTQUFTc0I7UUFDVHhCLFFBQVF5QjtJQUNWO0lBRUEsTUFBTUssU0FBUztRQUNicEMsd0RBQWMsQ0FBQztRQUNmUSxTQUFTO1FBQ1RGLFFBQVE7SUFDVjtJQUVBLHFCQUNFLDhEQUFDTCxZQUFZb0MsUUFBUTtRQUFDQyxPQUFPO1lBQUVqQztZQUFNRTtZQUFPRTtZQUFTaUI7WUFBT1E7WUFBVUU7WUFBUUcsaUJBQWlCLENBQUMsQ0FBQ2xDO1FBQUs7a0JBQ25HRDs7Ozs7O0FBR1AsRUFBQztBQUVNLE1BQU1vQyxVQUFVO0lBQ3JCLE1BQU1DLFVBQVU3QyxpREFBVUEsQ0FBQ0s7SUFDM0IsSUFBSSxDQUFDd0MsU0FBUztRQUNaLE1BQU0sSUFBSUMsTUFBTTtJQUNsQjtJQUNBLE9BQU9EO0FBQ1QsRUFBQyIsInNvdXJjZXMiOlsid2VicGFjazovL2Jlc2VuZC1mcm9udGVuZC8uL2NvbnRleHQvQXV0aENvbnRleHQudHN4P2ZkZmYiXSwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IFJlYWN0LCB7IGNyZWF0ZUNvbnRleHQsIHVzZUNvbnRleHQsIHVzZVN0YXRlLCB1c2VFZmZlY3QgfSBmcm9tICdyZWFjdCdcbmltcG9ydCBheGlvcyBmcm9tICdheGlvcydcbmltcG9ydCBDb29raWVzIGZyb20gJ2pzLWNvb2tpZSdcblxuaW50ZXJmYWNlIFVzZXIge1xuICBpZDogc3RyaW5nXG4gIGVtYWlsOiBzdHJpbmdcbiAgbmFtZTogc3RyaW5nXG4gIGNyZWF0ZWRBdDogc3RyaW5nXG59XG5cbmludGVyZmFjZSBBdXRoQ29udGV4dFR5cGUge1xuICB1c2VyOiBVc2VyIHwgbnVsbFxuICB0b2tlbjogc3RyaW5nIHwgbnVsbFxuICBsb2FkaW5nOiBib29sZWFuXG4gIGxvZ2luOiAoZW1haWw6IHN0cmluZywgcGFzc3dvcmQ6IHN0cmluZykgPT4gUHJvbWlzZTx2b2lkPlxuICByZWdpc3RlcjogKGVtYWlsOiBzdHJpbmcsIHBhc3N3b3JkOiBzdHJpbmcsIG5hbWU6IHN0cmluZykgPT4gUHJvbWlzZTx2b2lkPlxuICBsb2dvdXQ6ICgpID0+IHZvaWRcbiAgaXNBdXRoZW50aWNhdGVkOiBib29sZWFuXG59XG5cbmNvbnN0IEF1dGhDb250ZXh0ID0gY3JlYXRlQ29udGV4dDxBdXRoQ29udGV4dFR5cGUgfCB1bmRlZmluZWQ+KHVuZGVmaW5lZClcblxuZXhwb3J0IGNvbnN0IEF1dGhQcm92aWRlcjogUmVhY3QuRkM8eyBjaGlsZHJlbjogUmVhY3QuUmVhY3ROb2RlIH0+ID0gKHsgY2hpbGRyZW4gfSkgPT4ge1xuICBjb25zdCBbdXNlciwgc2V0VXNlcl0gPSB1c2VTdGF0ZTxVc2VyIHwgbnVsbD4obnVsbClcbiAgY29uc3QgW3Rva2VuLCBzZXRUb2tlbl0gPSB1c2VTdGF0ZTxzdHJpbmcgfCBudWxsPihudWxsKVxuICBjb25zdCBbbG9hZGluZywgc2V0TG9hZGluZ10gPSB1c2VTdGF0ZSh0cnVlKVxuXG4gIGNvbnN0IEFQSV9VUkwgPSBwcm9jZXNzLmVudi5ORVhUX1BVQkxJQ19BUElfVVJMIHx8ICdodHRwOi8vbG9jYWxob3N0OjgwODAnXG5cbiAgdXNlRWZmZWN0KCgpID0+IHtcbiAgICBjb25zdCBzdG9yZWRUb2tlbiA9IENvb2tpZXMuZ2V0KCdhdXRoX3Rva2VuJylcbiAgICBpZiAoc3RvcmVkVG9rZW4pIHtcbiAgICAgIHNldFRva2VuKHN0b3JlZFRva2VuKVxuICAgICAgZmV0Y2hVc2VyKHN0b3JlZFRva2VuKVxuICAgIH0gZWxzZSB7XG4gICAgICBzZXRMb2FkaW5nKGZhbHNlKVxuICAgIH1cbiAgfSwgW10pXG5cbiAgY29uc3QgZmV0Y2hVc2VyID0gYXN5bmMgKGF1dGhUb2tlbjogc3RyaW5nKSA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIGNvbnN0IHJlc3BvbnNlID0gYXdhaXQgYXhpb3MuZ2V0KGAke0FQSV9VUkx9L2FwaS91c2Vycy9tZWAsIHtcbiAgICAgICAgaGVhZGVyczogeyBBdXRob3JpemF0aW9uOiBgQmVhcmVyICR7YXV0aFRva2VufWAgfSxcbiAgICAgIH0pXG4gICAgICBzZXRVc2VyKHJlc3BvbnNlLmRhdGEpXG4gICAgfSBjYXRjaCAoZXJyb3IpIHtcbiAgICAgIGNvbnNvbGUuZXJyb3IoJ0ZhaWxlZCB0byBmZXRjaCB1c2VyOicsIGVycm9yKVxuICAgICAgQ29va2llcy5yZW1vdmUoJ2F1dGhfdG9rZW4nKVxuICAgICAgc2V0VG9rZW4obnVsbClcbiAgICB9IGZpbmFsbHkge1xuICAgICAgc2V0TG9hZGluZyhmYWxzZSlcbiAgICB9XG4gIH1cblxuICBjb25zdCBsb2dpbiA9IGFzeW5jIChlbWFpbDogc3RyaW5nLCBwYXNzd29yZDogc3RyaW5nKSA9PiB7XG4gICAgY29uc3QgcmVzcG9uc2UgPSBhd2FpdCBheGlvcy5wb3N0KGAke0FQSV9VUkx9L2FwaS9hdXRoL2xvZ2luYCwgeyBlbWFpbCwgcGFzc3dvcmQgfSlcbiAgICBjb25zdCB7IHRva2VuOiBuZXdUb2tlbiwgdXNlcjogbmV3VXNlciB9ID0gcmVzcG9uc2UuZGF0YVxuICAgIENvb2tpZXMuc2V0KCdhdXRoX3Rva2VuJywgbmV3VG9rZW4sIHsgZXhwaXJlczogNyB9KVxuICAgIHNldFRva2VuKG5ld1Rva2VuKVxuICAgIHNldFVzZXIobmV3VXNlcilcbiAgfVxuXG4gIGNvbnN0IHJlZ2lzdGVyID0gYXN5bmMgKGVtYWlsOiBzdHJpbmcsIHBhc3N3b3JkOiBzdHJpbmcsIG5hbWU6IHN0cmluZykgPT4ge1xuICAgIGNvbnN0IHJlc3BvbnNlID0gYXdhaXQgYXhpb3MucG9zdChgJHtBUElfVVJMfS9hcGkvYXV0aC9yZWdpc3RlcmAsIHsgZW1haWwsIHBhc3N3b3JkLCBuYW1lIH0pXG4gICAgY29uc3QgeyB0b2tlbjogbmV3VG9rZW4sIHVzZXI6IG5ld1VzZXIgfSA9IHJlc3BvbnNlLmRhdGFcbiAgICBDb29raWVzLnNldCgnYXV0aF90b2tlbicsIG5ld1Rva2VuLCB7IGV4cGlyZXM6IDcgfSlcbiAgICBzZXRUb2tlbihuZXdUb2tlbilcbiAgICBzZXRVc2VyKG5ld1VzZXIpXG4gIH1cblxuICBjb25zdCBsb2dvdXQgPSAoKSA9PiB7XG4gICAgQ29va2llcy5yZW1vdmUoJ2F1dGhfdG9rZW4nKVxuICAgIHNldFRva2VuKG51bGwpXG4gICAgc2V0VXNlcihudWxsKVxuICB9XG5cbiAgcmV0dXJuIChcbiAgICA8QXV0aENvbnRleHQuUHJvdmlkZXIgdmFsdWU9e3sgdXNlciwgdG9rZW4sIGxvYWRpbmcsIGxvZ2luLCByZWdpc3RlciwgbG9nb3V0LCBpc0F1dGhlbnRpY2F0ZWQ6ICEhdXNlciB9fT5cbiAgICAgIHtjaGlsZHJlbn1cbiAgICA8L0F1dGhDb250ZXh0LlByb3ZpZGVyPlxuICApXG59XG5cbmV4cG9ydCBjb25zdCB1c2VBdXRoID0gKCkgPT4ge1xuICBjb25zdCBjb250ZXh0ID0gdXNlQ29udGV4dChBdXRoQ29udGV4dClcbiAgaWYgKCFjb250ZXh0KSB7XG4gICAgdGhyb3cgbmV3IEVycm9yKCd1c2VBdXRoIG11c3QgYmUgdXNlZCB3aXRoaW4gQXV0aFByb3ZpZGVyJylcbiAgfVxuICByZXR1cm4gY29udGV4dFxufVxuIl0sIm5hbWVzIjpbIlJlYWN0IiwiY3JlYXRlQ29udGV4dCIsInVzZUNvbnRleHQiLCJ1c2VTdGF0ZSIsInVzZUVmZmVjdCIsImF4aW9zIiwiQ29va2llcyIsIkF1dGhDb250ZXh0IiwidW5kZWZpbmVkIiwiQXV0aFByb3ZpZGVyIiwiY2hpbGRyZW4iLCJ1c2VyIiwic2V0VXNlciIsInRva2VuIiwic2V0VG9rZW4iLCJsb2FkaW5nIiwic2V0TG9hZGluZyIsIkFQSV9VUkwiLCJwcm9jZXNzIiwiZW52IiwiTkVYVF9QVUJMSUNfQVBJX1VSTCIsInN0b3JlZFRva2VuIiwiZ2V0IiwiZmV0Y2hVc2VyIiwiYXV0aFRva2VuIiwicmVzcG9uc2UiLCJoZWFkZXJzIiwiQXV0aG9yaXphdGlvbiIsImRhdGEiLCJlcnJvciIsImNvbnNvbGUiLCJyZW1vdmUiLCJsb2dpbiIsImVtYWlsIiwicGFzc3dvcmQiLCJwb3N0IiwibmV3VG9rZW4iLCJuZXdVc2VyIiwic2V0IiwiZXhwaXJlcyIsInJlZ2lzdGVyIiwibmFtZSIsImxvZ291dCIsIlByb3ZpZGVyIiwidmFsdWUiLCJpc0F1dGhlbnRpY2F0ZWQiLCJ1c2VBdXRoIiwiY29udGV4dCIsIkVycm9yIl0sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./context/AuthContext.tsx\n");

/***/ }),

/***/ "./pages/_app.tsx":
/*!************************!*\
  !*** ./pages/_app.tsx ***!
  \************************/
/***/ ((module, __webpack_exports__, __webpack_require__) => {

"use strict";
eval("__webpack_require__.a(module, async (__webpack_handle_async_dependencies__, __webpack_async_result__) => { try {\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"default\": () => (/* binding */ App)\n/* harmony export */ });\n/* harmony import */ var react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react/jsx-dev-runtime */ \"react/jsx-dev-runtime\");\n/* harmony import */ var react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _context_AuthContext__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @/context/AuthContext */ \"./context/AuthContext.tsx\");\n/* harmony import */ var _styles_globals_css__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @/styles/globals.css */ \"./styles/globals.css\");\n/* harmony import */ var _styles_globals_css__WEBPACK_IMPORTED_MODULE_2___default = /*#__PURE__*/__webpack_require__.n(_styles_globals_css__WEBPACK_IMPORTED_MODULE_2__);\nvar __webpack_async_dependencies__ = __webpack_handle_async_dependencies__([_context_AuthContext__WEBPACK_IMPORTED_MODULE_1__]);\n_context_AuthContext__WEBPACK_IMPORTED_MODULE_1__ = (__webpack_async_dependencies__.then ? (await __webpack_async_dependencies__)() : __webpack_async_dependencies__)[0];\n\n\n\nfunction App({ Component, pageProps }) {\n    return /*#__PURE__*/ (0,react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxDEV)(_context_AuthContext__WEBPACK_IMPORTED_MODULE_1__.AuthProvider, {\n        children: /*#__PURE__*/ (0,react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxDEV)(Component, {\n            ...pageProps\n        }, void 0, false, {\n            fileName: \"/Users/bruno/infra/besend/frontend/pages/_app.tsx\",\n            lineNumber: 8,\n            columnNumber: 7\n        }, this)\n    }, void 0, false, {\n        fileName: \"/Users/bruno/infra/besend/frontend/pages/_app.tsx\",\n        lineNumber: 7,\n        columnNumber: 5\n    }, this);\n}\n\n__webpack_async_result__();\n} catch(e) { __webpack_async_result__(e); } });//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9wYWdlcy9fYXBwLnRzeCIsIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7O0FBQ29EO0FBQ3ZCO0FBRWQsU0FBU0MsSUFBSSxFQUFFQyxTQUFTLEVBQUVDLFNBQVMsRUFBWTtJQUM1RCxxQkFDRSw4REFBQ0gsOERBQVlBO2tCQUNYLDRFQUFDRTtZQUFXLEdBQUdDLFNBQVM7Ozs7Ozs7Ozs7O0FBRzlCIiwic291cmNlcyI6WyJ3ZWJwYWNrOi8vYmVzZW5kLWZyb250ZW5kLy4vcGFnZXMvX2FwcC50c3g/MmZiZSJdLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgdHlwZSB7IEFwcFByb3BzIH0gZnJvbSAnbmV4dC9hcHAnXG5pbXBvcnQgeyBBdXRoUHJvdmlkZXIgfSBmcm9tICdAL2NvbnRleHQvQXV0aENvbnRleHQnXG5pbXBvcnQgJ0Avc3R5bGVzL2dsb2JhbHMuY3NzJ1xuXG5leHBvcnQgZGVmYXVsdCBmdW5jdGlvbiBBcHAoeyBDb21wb25lbnQsIHBhZ2VQcm9wcyB9OiBBcHBQcm9wcykge1xuICByZXR1cm4gKFxuICAgIDxBdXRoUHJvdmlkZXI+XG4gICAgICA8Q29tcG9uZW50IHsuLi5wYWdlUHJvcHN9IC8+XG4gICAgPC9BdXRoUHJvdmlkZXI+XG4gIClcbn1cbiJdLCJuYW1lcyI6WyJBdXRoUHJvdmlkZXIiLCJBcHAiLCJDb21wb25lbnQiLCJwYWdlUHJvcHMiXSwic291cmNlUm9vdCI6IiJ9\n//# sourceURL=webpack-internal:///./pages/_app.tsx\n");

/***/ }),

/***/ "./styles/globals.css":
/*!****************************!*\
  !*** ./styles/globals.css ***!
  \****************************/
/***/ (() => {



/***/ }),

/***/ "react":
/*!************************!*\
  !*** external "react" ***!
  \************************/
/***/ ((module) => {

"use strict";
module.exports = require("react");

/***/ }),

/***/ "react/jsx-dev-runtime":
/*!****************************************!*\
  !*** external "react/jsx-dev-runtime" ***!
  \****************************************/
/***/ ((module) => {

"use strict";
module.exports = require("react/jsx-dev-runtime");

/***/ }),

/***/ "axios":
/*!************************!*\
  !*** external "axios" ***!
  \************************/
/***/ ((module) => {

"use strict";
module.exports = import("axios");;

/***/ }),

/***/ "js-cookie":
/*!****************************!*\
  !*** external "js-cookie" ***!
  \****************************/
/***/ ((module) => {

"use strict";
module.exports = import("js-cookie");;

/***/ })

};
;

// load runtime
var __webpack_require__ = require("../webpack-runtime.js");
__webpack_require__.C(exports);
var __webpack_exec__ = (moduleId) => (__webpack_require__(__webpack_require__.s = moduleId))
var __webpack_exports__ = (__webpack_exec__("./pages/_app.tsx"));
module.exports = __webpack_exports__;

})();