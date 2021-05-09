use tokio::net::UnixDatagram;
use std::io;

#[tokio::main]
async fn main() -> io::Result<()> {
    let dir = tempfile::tempdir().unwrap();
    let socket_path = dir.path().join("git_trace.sock");

    let mut socket = UnixDatagram::bind(&socket_path)?;
    socket.connect(&socket_path)?;
    println!("Connected to {:?}", &socket_path);

    let mut data = vec![0; 1024];
    loop {
        match socket.try_recv(&mut data[..]) {
            Ok(size) => {
                let dgram = &data[..size];
                let event = String::from_utf8(dgram.to_vec()).unwrap();
                handle_event(event).await;
            }

            // False-positive, continue
            Err(ref e) if e.kind() == io::ErrorKind::WouldBlock => {}
            Err(e) => {
                return Err(e);
            }
        }
    }
}

async fn handle_event(event: String) {
    println!("{}", event);
}
