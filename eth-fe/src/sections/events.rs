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
pub struct EthFormEventAddress {
    pub eth_event_address: String,
}

fn format_json(input: &str) -> Result<String, Error> {
    let value: Value = serde_json::from_str(input)?;
    serde_json::to_string_pretty(&value)
}

#[function_component(Events)]
pub fn events() -> Html {
    let eth_event_address_ref = use_node_ref();
    let eth_event_address_error = use_state(|| String::new());
    let is_loading_event_address = use_state(|| false);
    let response_data_event_address = use_state(|| String::new());

    let onsubmit_event_address = {
        let eth_event_address_ref = eth_event_address_ref.clone();
        let eth_event_address_error = eth_event_address_error.clone();
        let is_loading_event_address = is_loading_event_address.clone();
        let response_data_event_address = response_data_event_address.clone();

        Callback::from(move |event: SubmitEvent| {
            event.prevent_default();

            response_data_event_address.set(String::new());
            eth_event_address_error.set(String::new());

            let eth_event_address = eth_event_address_ref.cast::<HtmlInputElement>().unwrap().value();

            if eth_event_address.is_empty() {
                eth_event_address_error.set("Event address cannot be empty.".to_string());
                return;
            }

            if eth_event_address.len() < 40 {
                eth_event_address_error.set("Event address is too short.".to_string());
                return;
            }

            if eth_event_address.len() > 50 {
                eth_event_address_error.set("Event address is too long.".to_string());
                return;
            }

            if !Regex::new(r"^(0x)?[a-fA-F0-9]{40}$").unwrap().is_match(&eth_event_address) {
                eth_event_address_error.set("Event address does not match the required pattern.".to_string());
                return;
            }

            let eth_event_address_form = EthFormEventAddress {
                eth_event_address,
            };


            wasm_bindgen_futures::spawn_local({
                let is_loading_event_address = is_loading_event_address.clone();
                let response_data_event_address = response_data_event_address.clone();

                async move {
                    is_loading_event_address.set(true);

                    let url = format!("{}eth-events/{}", *BASE_URL, &eth_event_address_form.eth_event_address);

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

                    is_loading_event_address.set(false);
                    response_data_event_address.set(formatted_json);
                }
            });
        })
    };


    html! {
        <div>
            <form onsubmit={onsubmit_event_address.clone()}>
                <InputField input_node_ref={eth_event_address_ref.clone()} label={"Event Address".to_owned()}  name={"eth_event_address".clone()} field_type={"text".clone()} loading={*is_loading_event_address.clone()}  />
                <span class="text-danger">{ &*eth_event_address_error.clone() }</span>
            </form>
            {
                if *is_loading_event_address.clone() {
                    html! { <div class="spinner-border" role="status"><span class="sr-only"></span></div> }
                } else {
                    html! {}
                }
            }
            <pre class="bg-dark p-3 border" style={if response_data_event_address.is_empty() { "display: none;" } else { "display: block;" }}>
                <code class="language-html">{ &*response_data_event_address.clone() }</code>
            </pre>
        </div>
    }
}