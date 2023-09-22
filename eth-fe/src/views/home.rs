use yew::prelude::*;
use crate::sections::blocks::Blocks;
use crate::sections::transactions::Transactions;
use crate::sections::events::Events;
use crate::sections::latests::Latests;

#[function_component(Home)]
pub fn home() -> Html {
    let active_tab = use_state(|| String::from("block"));
    let set_active_tab = Callback::from({
        let active_tab = active_tab.clone();
        move |tab: String| active_tab.set(tab)
    });

    html! {
        <main class="container">
            <div class="px-4 py-5 my-5 text-center">
                <div class="eth-img"></div>
                <h1 class="display-5 fw-bold text-body-emphasis">{"Ethereum Block Service"}</h1>
            </div>
            <div class="container-fluid">
                <div class="row">
                    <div class="col">
                        <nav>
                            <div class="nav nav-tabs" id="nav-tab" role="tablist">
                                <button
                                    class={if *active_tab == "block" { "nav-link active" } else { "nav-link" }}
                                    id="nav-block-tab"
                                    type="button"
                                    role="tab"
                                    onclick={Callback::from({
                                    let set_active_tab = set_active_tab.clone();
                                    move |_| set_active_tab.emit(String::from("block"))
                                })}>
                                        {"Block"}
                                </button>
                                <button
                                    class={if *active_tab == "transaction" { "nav-link active" } else { "nav-link" }}
                                    id="nav-block-tab"
                                    type="button"
                                    role="tab"
                                    onclick={Callback::from({
                                    let set_active_tab = set_active_tab.clone();
                                    move |_| set_active_tab.emit(String::from("transaction"))
                                })}>
                                        {"Transaction"}
                                </button>
                                <button
                                    class={if *active_tab == "event" { "nav-link active" } else { "nav-link" }}
                                    id="nav-block-tab"
                                    type="button"
                                    role="tab"
                                    onclick={Callback::from({
                                    let set_active_tab = set_active_tab.clone();
                                    move |_| set_active_tab.emit(String::from("event"))
                                })}>
                                        {"Event"}
                                </button>
                                <button
                                    class={if *active_tab == "latest" { "nav-link active" } else { "nav-link" }}
                                    id="nav-block-tab"
                                    type="button"
                                    role="tab"
                                    onclick={Callback::from({
                                    let set_active_tab = set_active_tab.clone();
                                    move |_| set_active_tab.emit(String::from("latest"))
                                })}>
                                        {"Latest Blocks"}
                                </button>
                            </div>
                        </nav>
                        <div class="tab-content" id="nav-tabContent">
                            <div
                                class={if *active_tab == "block" { "tab-pane fade show active" } else { "tab-pane fade" }}
                                id="nav-block"
                                role="tabpanel"
                                tabindex="0">
                                    <Blocks />
                            </div>
                            <div
                                class={if *active_tab == "transaction" { "tab-pane fade show active" } else { "tab-pane fade" }}
                                id="nav-transaction"
                                role="tabpanel"
                                tabindex="1">
                                    <Transactions />
                            </div>
                            <div
                                class={if *active_tab == "event" { "tab-pane fade show active" } else { "tab-pane fade" }}
                                id="nav-event"
                                role="tabpanel"
                                tabindex="2">
                                    <Events />
                            </div>
                            <div
                                class={if *active_tab == "latest" { "tab-pane fade show active" } else { "tab-pane fade" }}
                                id="nav-latest"
                                role="tabpanel"
                                tabindex="3">
                                    <Latests />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </main>
    }
}