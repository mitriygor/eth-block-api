use gloo_net::http::{Headers, Request};
use once_cell::sync::Lazy;
use serde::{Deserialize, Serialize};
use web_sys::HtmlInputElement;
use serde_json::{Value, Error};

use yew::prelude::*;

use crate::components::input::InputField;

static BASE_URL: Lazy<String> = Lazy::new(|| {
    std::env::var("BASE_URL").unwrap_or_else(|_| "http://localhost:3000/".to_string())
});

#[derive(Clone, PartialEq, Properties, Debug, Default, Serialize, Deserialize)]
pub struct EthFormBlockId {
    pub eth_block_id: String,
}

fn format_json(input: &str) -> Result<String, Error> {
    let value: Value = serde_json::from_str(input)?;
    serde_json::to_string_pretty(&value)
}

#[function_component(Blocks)]
pub fn blocks() -> Html {
    let eth_block_id_ref = use_node_ref();
    let eth_block_id_error = use_state(|| String::new());
    let is_loading_block_id = use_state(|| false);
    let response_data_block_id = use_state(|| String::new());

    let onsubmit_block_id = {
        let eth_block_id_ref = eth_block_id_ref.clone();
        let eth_block_id_error = eth_block_id_error.clone();
        let is_loading_block_id = is_loading_block_id.clone();
        let response_data_block_id = response_data_block_id.clone();

        Callback::from(move |event: SubmitEvent| {
            event.prevent_default();

            response_data_block_id.set(String::new());
            eth_block_id_error.set(String::new());

            let eth_block_id = eth_block_id_ref.cast::<HtmlInputElement>().unwrap().value();

            if eth_block_id.is_empty() {
                eth_block_id_error.set("ID cannot be empty.".to_string());
                return;
            }
            if eth_block_id.len() < 5 {
                eth_block_id_error.set("ID is too short.".to_string());
                return;
            }
            if eth_block_id.len() > 70 {
                eth_block_id_error.set("ID is too long.".to_string());
                return;
            }

            let eth_block_id_form = EthFormBlockId {
                eth_block_id,
            };

            wasm_bindgen_futures::spawn_local({
                let is_loading_block_id = is_loading_block_id.clone(); // Clone again here
                let response_data_block_id = response_data_block_id.clone();

                async move {
                    is_loading_block_id.set(true);

                    let url = format!("{}eth-blocks/{}", *BASE_URL, &eth_block_id_form.eth_block_id);

                    let get_request = Request::get(&url)
                        .headers({
                            let headers = Headers::new();
                            headers.set("Content-Type", "application/json");
                            headers
                        })
                        .send()
                        .await
                        .expect("Failed to send the request");

                    let response_text = get_request.text().await.expect("{'error': 'Failed to read the response text'}");
                    let formatted_json = format_json(&response_text).unwrap_or_else(|_| response_text.clone());

                    is_loading_block_id.set(false);
                    response_data_block_id.set(formatted_json);
                }
            });
        })
    };


    html! {
        <div>
            <form onsubmit={onsubmit_block_id}>
                <InputField input_node_ref={eth_block_id_ref} label={"Block Identifier:".to_owned()}  name={"eth_block_id".clone()} field_type={"text".clone()} loading={*is_loading_block_id.clone()}  />
                <span class="text-danger">{ &*eth_block_id_error.clone() }</span>
            </form>
            {
                if *is_loading_block_id.clone() {
                    html! { <div class="spinner-border" role="status"><span class="sr-only"></span></div> }
                } else {
                    html! {}
                }
            }
            <pre class="bg-dark p-3 border" style={if response_data_block_id.is_empty() { "display: none;" } else { "display: block;" }}>
                <code class="language-html">{ &*response_data_block_id.clone() }</code>
            </pre>
        </div>
    }
}