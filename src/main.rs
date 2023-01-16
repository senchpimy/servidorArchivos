use rocket::response::Redirect;
#[macro_use] extern crate rocket;

#[get("/<_..>")]
fn everything() -> Redirect{
        Redirect::to("/")
}

#[get("/")]
fn papu()-> String{
"hola mundo".to_string()
}

#[launch]
fn rocket() -> _ {
    rocket::build().mount("/", routes![papu]).mount("/", routes![everything])
    //rocket::build().mount("/", routes![hello])
}

