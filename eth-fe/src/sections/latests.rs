use gloo_net::http::{Headers, Request};
use once_cell::sync::Lazy;
use serde_json::{Value, Error};


use yew::prelude::*;

static BASE_URL: Lazy<String> = Lazy::new(|| {
    std::env::var("BASE_URL").unwrap_or_else(|_| "http://localhost:3000/".to_string())
});

fn format_json(input: &str) -> Result<String, Error> {
    let value: Value = serde_json::from_str(input)?;
    serde_json::to_string_pretty(&value)
}

#[function_component(Latests)]
pub fn latests() -> Html {
    let is_loading_latests = use_state(|| false);
    let response_data_latests = use_state(|| String::new());

    let onfetch_latests = {
        let is_loading_latests = is_loading_latests.clone();
        let response_data_latests = response_data_latests.clone();

        Callback::from(move |event: MouseEvent| {
            event.prevent_default();

            response_data_latests.set(String::new());

            wasm_bindgen_futures::spawn_local({
                let is_loading_latests = is_loading_latests.clone();
                let response_data_latests = response_data_latests.clone();

                async move {
                    is_loading_latests.set(true);

                    let url = format!("{}eth-blocks/latest", *BASE_URL);

                    let get_request = Request::get(&url)
                        .headers({
                            let headers = Headers::new();
                            headers.set("Content-Type", "application/json");
                            headers
                        })
                        .send()
                        .await
                        .expect("Failed to send the request");

                    let response_text = get_request.text().await.expect("Failed to read the response text");
                    let formatted_json = format_json(&response_text).unwrap_or_else(|_| response_text.clone());

                    is_loading_latests.set(false);
                    response_data_latests.set(formatted_json);
                }
            });
        })
    };


    html! {
        <div>
            <button type="button" class="btn btn-outline-secondary refresh-button" onclick={onfetch_latests}><i class="bi bi-arrow-clockwise"></i></button>
            {
                if *is_loading_latests.clone() {
                    html! { <div class="spinner-border" role="status"><span class="sr-only"></span></div> }
                } else {
                    html! {}
                }
            }
            <pre class="bg-dark p-3 border" style={if response_data_latests.is_empty() { "display: none;" } else { "display: block;" }}>
                <code class="language-html">{ &*response_data_latests.clone() }</code>
            </pre>
        </div>
    }
}