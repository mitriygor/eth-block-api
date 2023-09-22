use gloo_net::http::{Headers, Request};
use once_cell::sync::Lazy;
use serde::{Deserialize, Serialize};
use web_sys::HtmlInputElement;
use serde_json::{Value, Error};

use yew::prelude::*;
use regex::Regex;

use crate::components::input::InputField;

static BASE_URL: Lazy<String> = Lazy::new(|| {
    std::env::var("BASE_URL").unwrap_or_else(|_| "http://localhost:3000/".to_string())
});

#[derive(Clone, PartialEq, Properties, Debug, Default, Serialize, Deserialize)]
pub struct EthFormTransactionHash {
    pub eth_transaction_hash: String,
}

fn format_json(input: &str) -> Result<String, Error> {
    let value: Value = serde_json::from_str(input)?;
    serde_json::to_string_pretty(&value)
}

#[function_component(Transactions)]
pub fn transactions() -> Html {
    let eth_transaction_hash_ref = use_node_ref();
    let eth_transaction_hash_error = use_state(|| String::new());
    let is_loading_transaction_hash = use_state(|| false);
    let response_data_transaction_hash = use_state(|| String::new());

    let onsubmit_transaction_hash = {
        let eth_transaction_hash_ref = eth_transaction_hash_ref.clone();
        let eth_transaction_hash_error = eth_transaction_hash_error.clone();
        let is_loading_transaction_hash = is_loading_transaction_hash.clone();
        let response_data_transaction_hash = response_data_transaction_hash.clone();

        Callback::from(move |event: SubmitEvent| {
            event.prevent_default();

            response_data_transaction_hash.set(String::new());
            eth_transaction_hash_error.set(String::new());

            let eth_transaction_hash = eth_transaction_hash_ref.cast::<HtmlInputElement>().unwrap().value();

            if eth_transaction_hash.is_empty() {
                eth_transaction_hash_error.set("Transaction hash cannot be empty.".to_string());
                return;
            }

            if eth_transaction_hash.len() < 60 {
                eth_transaction_hash_error.set("Transaction hash is too short.".to_string());
                return;
            }

            if eth_transaction_hash.len() > 70 {
                eth_transaction_hash_error.set("Transaction hash is too long.".to_string());
                return;
            }

            if !Regex::new(r"^(0x)?[a-fA-F0-9]{1,64}$").unwrap().is_match(&eth_transaction_hash) {
                eth_transaction_hash_error.set("Transaction hash does not match the required pattern.".to_string());
                return;
            }

            let eth_transaction_hash_form = EthFormTransactionHash {
                eth_transaction_hash,
            };

            log::info!("eth_transaction_hash_form {:?}", &eth_transaction_hash_form);

            wasm_bindgen_futures::spawn_local({
                let is_loading_transaction_hash = is_loading_transaction_hash.clone();
                let response_data_transaction_hash = response_data_transaction_hash.clone();

                async move {
                    is_loading_transaction_hash.set(true);

                    let url = format!("{}eth-transactions/{}", *BASE_URL, &eth_transaction_hash_form.eth_transaction_hash);

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

                    is_loading_transaction_hash.set(false);
                    response_data_transaction_hash.set(formatted_json);
                }
            });
        })
    };


    html! {
        <div>
            <form onsubmit={onsubmit_transaction_hash.clone()}>
                <InputField input_node_ref={eth_transaction_hash_ref.clone()} label={"Transaction Hash".to_owned()}  name={"eth_transaction_hash".clone()} field_type={"text".clone()} loading={*is_loading_transaction_hash.clone()}  />
                <span class="text-danger">{ &*eth_transaction_hash_error.clone() }</span>
            </form>
            {
                if *is_loading_transaction_hash.clone() {
                    html! { <div class="spinner-border" role="status"><span class="sr-only"></span></div> }
                } else {
                    html! {}
                }
            }
            <pre class="bg-dark p-3 border" style={if response_data_transaction_hash.is_empty() { "display: none;" } else { "display: block;" }}>
                <code class="language-html">{ &*response_data_transaction_hash.clone() }</code>
            </pre>
        </div>
    }
}