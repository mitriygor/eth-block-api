use gloo_net::http::{Headers, Request};
use serde::{Deserialize, Serialize};
use serde_json;
use wasm_bindgen::JsValue;
use web_sys::HtmlInputElement;

use yew::prelude::*;

use crate::components::input::InputField;

#[derive(Clone, PartialEq, Properties, Debug, Default, Serialize, Deserialize)]
pub struct EthForm {
    pub eth_block_id: String,
    pub eth_transaction_hash: String,
    pub eth_event_address: String,
}

#[function_component(Home)]
pub fn home() -> Html {
    let eth_form = use_state(|| EthForm::default());

    let eth_block_id = use_node_ref();
    let eth_transaction_hash_ref = use_node_ref();
    let eth_event_address_ref = use_node_ref();

    log::info!("eth_form {:?}", &eth_form.clone());
    let onsubmit = {
        let eth_form = eth_form.clone();

        let eth_block_id = eth_block_id.clone();
        let eth_transaction_hash_ref = eth_transaction_hash_ref.clone();
        let eth_event_address_ref = eth_event_address_ref.clone();

        Callback::from(move |event: SubmitEvent| {
            event.prevent_default();
            log::info!("eth_form {:?}", &eth_form.clone());

            let eth_block_id = eth_block_id.cast::<HtmlInputElement>().unwrap().value();
            let eth_transaction_hash = eth_transaction_hash_ref.cast::<HtmlInputElement>().unwrap().value();
            let eth_event_address = eth_event_address_ref.cast::<HtmlInputElement>().unwrap().value();

            let eth_form = EthForm {
                eth_block_id,
                eth_transaction_hash,
                eth_event_address,
            };

            log::info!("eth_form {:?}", &eth_form);

            wasm_bindgen_futures::spawn_local(async move {
                let post_request = Request::post("#")
                    .headers({
                        let headers = Headers::new();
                        headers
                            // .append(name, value)
                            .set("Content-Type", "application/json");
                        headers
                    })
                    .body(JsValue::from(
                        serde_json::to_string(&eth_form).unwrap(),
                    ))
                    .send()
                    .await
                    .unwrap();

                log::info!("post_request {:?}", &post_request);
            });
        })
    };

    html! {
        <main class="home">
            <form {onsubmit} class="eth-form">
                <InputField input_node_ref={eth_block_id} label={"Block Identifier".to_owned()} name={"eth_block_id".clone()} field_type={"text".clone()} />
                <InputField input_node_ref={eth_transaction_hash_ref} label={"Transaction Hash".to_owned()} name={"eth_transaction_hash".clone()} field_type={"text".clone()}  />
                <InputField input_node_ref={eth_event_address_ref} label={"Event Address".to_owned()} name={"eth_event_address".clone()} field_type={"text".clone()}  />
                <button type="submit" class="button button-primary">{"Submit"}</button>
            </form>
        </main>
    }
}