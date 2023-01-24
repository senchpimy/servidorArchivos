use rocket::response::Redirect;
use rocket::config::Config;
use rocket::fs::NamedFile;
use json;
use std::path::Path;
use std::fs;
#[macro_use] extern crate rocket;

const FILES_TO_SERVE: &str= "/home/plof/Documents/RustProjects/servidorArchivos/config.json";
const TRIES_TO_GET_A_FILE: i32 = 10;

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
    let data= json::parse(input).unwrap();
   if data["name"]=="neco" {
        let file = include_bytes!("../neco-arc.gif");
        print!("A wild Neco-arc appeared");
        return file
    }else{
        let contents = fs::read_to_string(&FILES_TO_SERVE)
            .expect("Should have been able to read the file");
        let parsed = json::parse(&contents).unwrap();
        for (key, _) in parsed.entries() {
            if parsed[key]==data["key"]{
                if parsed[key]["password"]=="none"{
                    return serve_file(&data["file"], &parsed[key]["files"]);
                }else{
                    if parsed[key]["password"]==data["password"]{
                        return serve_file(&data["file"], &parsed[key]["files"]);
                    }else{
                        println!("Wrong password");
                        return br"Wrong password";
                        //Add one to attemps getting the file block at specified CONST
                    }
                }
            }
            println!("The name of the share wasnt found");
            return br"The name of the share wasnt found";
        }
    }
    println!("Server Error");
    br"Server Error"
}

#[launch]
fn rocket() -> _ {
    let config = Config::default();
    rocket::custom(config).mount("/", routes![index]).mount("/", routes![everything])
                          .mount("/",routes![request_handler])
}

fn serve_file(path:&json::JsonValue,files:&json::JsonValue)->Vec<u8>{
    if files.is_array()==false{
        println!("Error in the file configuration, files must be an array");
        return b"Error in the file configuration, files must be an array".to_vec()
    }
    if path.is_string()==false{
        println!("Error in the request, requested path must be a string");
        return b"Error in the request, requested path must be a string".to_vec()
    }
    if Path::new(path.as_str().unwrap()).exists()==false{
        println!("File requested does not exist");
        return b"File requested does not exist".to_vec()
    }
    for i in 0..files.len(){
        if files[i].is_string()==false{
            println!("Error in configuration, path is not a string");
            return b"Error in configuration, path is not a string".to_vec()
        }
        if files[i].as_str()==path.as_str(){
    let file = fs::read(Path::new(files[i].as_str().unwrap())).unwrap();
    return file;
        }
        println!("Cannot serve file");
        return b"Cannot serve file".to_vec()
    }
    println!("Error sharing the file");
    b"Error sharing the file".to_vec()
}

//fn serve_file(path:&json::JsonValue,files:&json::JsonValue)->&'static [u8]{
//    if files.is_array()==false{
//        println!("Error in the file configuration, files must be an array");
//        return b"Error in the file configuration, files must be an array"
//    }
//    if path.is_string()==false{
//        println!("Error in the request, requested path must be a string");
//        return b"Error in the request, requested path must be a string"
//    }
//    if Path::new(path.as_str().unwrap()).exists()==false{
//        println!("File requested does not exist");
//        return b"File requested does not exist"
//    }
//    for i in 0..files.len(){
//        if files[i].is_string()==false{
//            println!("Error in configuration, path is not a string");
//            return b"Error in configuration, path is not a string"
//        }
//        if files[i].as_str()==path.as_str(){
//            let file:&[u8] = &fs::read(Path::new(files[i].as_str().unwrap())).unwrap();
//            return file;
//        }
//        println!("Cannot serve file");
//        return b"Cannot serve file"
//    }
//    println!("Error sharing the file");
//    b"Error sharing the file"
//}
