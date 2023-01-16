use rocket::response::Redirect;
use rocket::config::Config;
use rocket::fs::NamedFile;
use rocket::serde::{Deserialize, json};
#[macro_use] extern crate rocket;

#[derive(Debug, PartialEq, Deserialize)]
#[serde(crate = "rocket::serde")]
struct Data<'r> {
    password: &'r str,
    file: &'r str,
}

#[get("/<_..>")]
fn everything() -> Redirect{
        Redirect::to("/")
}

#[get("/")]
async fn index() -> Option<NamedFile> {
    NamedFile::open("index.html").await.ok()
}

#[post("/",data="<input>")]
fn request_handler(input:&str)->&'static [u8]{
    println!("papu");
    println!("{}",input);
    let data: Data = json::from_str(input).unwrap();
   if data.file=="neco" {
        let file = include_bytes!("../neco-arc.gif");
        return file
    }
    br"No hola"
}

#[launch]
fn rocket() -> _ {
    let config = Config::default();
    rocket::custom(config).mount("/", routes![index]).mount("/", routes![everything])
                          .mount("/",routes![request_handler])
}
