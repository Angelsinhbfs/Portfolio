import React, { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import App from "./App";
import {Auth0Provider} from "@auth0/auth0-react";

const rootElement = document.getElementById("root");
if (rootElement) {
    const root = createRoot(rootElement);
    root.render(
        <StrictMode>
            <Auth0Provider
                domain="dev-kokgswdod5i61zql.us.auth0.com"
                clientId="y3vpo3ZXk4GfiQzjMOPtuxazwNxaDjLL"
                authorizationParams={{
                    redirect_uri: window.location.origin
                }}
            >
                <App />
            </Auth0Provider>
        </StrictMode>
    );
} else {
    console.error("Root element not found in the DOM");
}