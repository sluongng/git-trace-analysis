use std::io;
use tokio::net::UnixDatagram;

mod event;
pub use self::event::Event;

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
                handle_event(dgram).await;
            }

            // False-positive, continue
            Err(ref e) if e.kind() == io::ErrorKind::WouldBlock => {}
            Err(e) => {
                return Err(e);
            }
        }
    }
}

async fn handle_event(event: &[u8]) {
    let e: Event = match serde_json::from_slice(event) {
        Err(_e) => return,
        Ok(v) => v,
    };

    println!("{}", e.sid);
}

#[cfg(test)]
mod tests {
    use super::*;
    use tokio_test::block_on;

    #[test]
    fn should_handle_event() {
        block_on(
            handle_event(
                br#"{"event":"atexit","sid":"20210509T164222.307807Z-Hf8bd5587-P00013e6d","thread":"main","time":"2021-05-09T16:42:31.768234Z","file":"trace2/tr2_tgt_event.c","line":201,"t_abs":9.461512,"code":0}"#
            )
        );
    }
}
