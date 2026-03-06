use axum::{Router, response::Html, routing::get};
use tokio::fs;
use tokio::net::TcpListener;

#[tokio::main]
async fn main() {
    let app_route = Router::new()
        .route("/", get(main_page_handler))
        .route("/about-us", get(about_page_handler))
        .route("/buy-merch", get(buy_merch_page_handler));

    let listener = TcpListener::bind("0.0.0.0:8080").await.unwrap();

    axum::serve(listener, app_route).await.unwrap();
}

async fn main_page_handler() -> Html<String> {
    let html = fs::read_to_string("web/templates/index.html")
        .await
        .unwrap();
    Html(html)
}

async fn about_page_handler() -> Html<String> {
    let html = fs::read_to_string("web/templates/about.html")
        .await
        .unwrap();
    Html(html)
}

async fn buy_merch_page_handler() -> Html<String> {
    let html = fs::read_to_string("web/templates/buyMerch.html")
        .await
        .unwrap();
    Html(html)
}
