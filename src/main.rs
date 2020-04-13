use async_std::task;
use async_std::prelude::*;
use async_std::os::unix::net::UnixListener;
use async_std::path::Path;
use async_std::fs::remove_file;

mod event;

async fn handle_events(events: String) {
    events.lines().for_each(|l| {
        let event: event::Event = serde_json::from_str(l).unwrap();

        println!("{}", event.event);
    });
}

#[async_std::main]
async fn main() -> std::io::Result<()> {
    let sock_address = "/tmp/git_trace.sock";

    // Check if socket already exists.
    let socket_path = Path::new(sock_address);
    if socket_path.exists().await {
        println!("Removing old socket {}", sock_address);
        remove_file(socket_path).await?
    }

    // Open new socket listener
    let listener = UnixListener::bind(sock_address).await?;
    println!("Listening on {}", sock_address);

    let mut incoming = listener.incoming();

    while let Some(stream) = incoming.next().await {
        let mut events = String::new();
        stream?
            .read_to_string(&mut events)
            .await?;

        task::spawn(handle_events(events));
    }

    Ok(())
}
