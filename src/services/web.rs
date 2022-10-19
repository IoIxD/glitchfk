use std::net::SocketAddr;
use hyper::{Body, Request, Response, Server};
use hyper::service::Service;

use std::future::Future;
use std::pin::Pin;
use std::task::{Context, Poll};

use crate::modules::gradient;
use crate::image;

pub async fn web_thread() -> Result<&'static str, hyper::Error> {   
    let addr = SocketAddr::from(([127,0,0,1], 6969));

    // Then bind and serve...
    let server = Server::bind(&addr).serve(MakeSvc{});

    // And run forever...
    match server.await {
        Ok(_) => Ok(""),
        Err(e) => Err(e),
    }
}

struct Svc {
}

impl Service<Request<Body>> for Svc {
    type Response = Response<Body>;
    type Error = hyper::Error;
    
    type Future = Pin<Box<dyn Future<Output = Result<Self::Response, Self::Error>> + Send>>;

    fn poll_ready(&mut self, _: &mut Context) -> Poll<Result<(), Self::Error>> {
        Poll::Ready(Ok(()))
    }

    fn call(&mut self, _req: Request<Body>) -> Self::Future {
        // TODO: query values
        let grad1 = gradient::random_gradient();
        let grad2 = gradient::random_gradient();
        let final_grad = image::xor_images(grad1, grad2);

        let mut pixels_raw: Vec<u8> = vec![0; 0];
        // TODO: can we make this a one liner?
        for p in final_grad.pixels() {
            pixels_raw.push(p.0[0]);
            pixels_raw.push(p.0[1]);
            pixels_raw.push(p.0[2]);
        };

        let pixels = image::png_from_u8(pixels_raw);
        let res = Ok(Response::builder().body(Body::from(pixels)).unwrap());
        Box::pin(async { res })


    }
}

struct MakeSvc {
}

impl<T> Service<T> for MakeSvc {
    type Response = Svc;
    type Error = hyper::Error;
    type Future = Pin<Box<dyn Future<Output = Result<Self::Response, Self::Error>> + Send>>;

    fn poll_ready(&mut self, _: &mut Context) -> Poll<Result<(), Self::Error>> {
        Poll::Ready(Ok(()))
    }

    fn call(&mut self, _: T) -> Self::Future {
        let fut = async move { Ok(Svc {}) };
        Box::pin(fut)
    }
}