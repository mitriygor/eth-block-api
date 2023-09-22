mod app;
mod components;
mod views;
mod sections;

use app::App;
use dotenv::dotenv;

fn main() {
    dotenv().ok();
    wasm_logger::init(wasm_logger::Config::default());
    yew::Renderer::<App>::new().render();
}