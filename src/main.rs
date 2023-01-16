use rocket::response::Redirect;
use rocket::fs::NamedFile;
#[macro_use] extern crate rocket;

#[get("/<_..>")]
fn everything() -> Redirect{
        Redirect::to("/")
}

#[get("/")]
async fn index() -> Option<NamedFile> {
    NamedFile::open("index.html").await.ok()
}

#[launch]
fn rocket() -> _ {
    rocket::build().mount("/", routes![index]).mount("/", routes![everything])
}
