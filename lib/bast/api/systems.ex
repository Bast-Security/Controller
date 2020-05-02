defmodule Bast.Api.Systems do
  use Plug.Router

  alias Bast.Api.Status

  plug :match
  plug :dispatch

  post "/" do
    Status.send_status(Status.not_implemented, conn)
  end

  delete "/" do
    Status.send_status(Status.not_implemented, conn)
  end

  get "/" do
    Status.send_status(Status.not_implemented, conn)
  end

  get "/totp" do
    Status.send_status(Status.not_implemented, conn)
  end

  get "/log" do
    Status.send_status(Status.not_implemented, conn)
  end

  forward "/users", to: Bast.Api.Systems.Users
  forward "/locks", to: Bast.Api.Systems.Locks
  forward "/roles", to: Bast.Api.Systems.Roles

  match _ do
    Status.send_status(Status.not_found, conn)
  end
end
