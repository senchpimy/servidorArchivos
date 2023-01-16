use rocket::response::Redirect;
use rocket::config::Config;
use rocket::fs::NamedFile;
use rocket::serde::{Deserialize, json};
#[macro_use] extern crate rocket;

#[derive(Debug, PartialEq, Deserialize)]
#[serde(crate = "rocket::serde")]
struct Data<'r> {
    ip: &'r str,
    password: &'r str,
    stars: usize,
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
fn test(input:&str)->&'static str{
    println!("papu");
    println!("{}",input);
let data: Data = json::from_str(input).unwrap();
if (data==Data { ip: "Rocket",password: "Rocket", stars: 5, }){
    return "hola"
}
"No hola"
}

#[launch]
fn rocket() -> _ {
    let config = Config::default();
    rocket::custom(config).mount("/", routes![index]).mount("/", routes![everything])
                          .mount("/",routes![test])
}
