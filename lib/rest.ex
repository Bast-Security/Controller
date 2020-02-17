defmodule Bast.Controller.REST do
  @moduledoc """
  This is the HTTP server responsible for handling HTTP
  REST messages sent to the middleman.
  """

  use Plug.Router

  plug Plug.Logger
  plug Plug.Parsers,
    parsers: [:json,:urlencoded,:multipart],
    json_decoder: Jason
  plug :match
  plug :dispatch

  post "/addUser" do
    body = conn.body_params

    user = %Bast.Controller.User{
      name: body["name"],
      email: body["email"],
      pin: body["pin"],
      cardno: body["cardno"]
    } |> Bast.Controller.User.validate

    case Bast.Controller.Repo.insert(user, on_conflict: :nothing) do
      {:ok, _user} -> send_resp(conn, 200, "User added!\n")
      {:error, _changeset} -> send_resp(conn, 400, "Malformed request!\n")
    end
  end

  post "/addRole" do
    user = %Bast.Controller.Role{
      name: conn.body_params["name"]
    } |> Bast.Controller.Role.validate

    case Bast.Controller.Repo.insert(user, on_conflict: :nothing) do
      {:ok, _user} -> send_resp(conn, 200, "Role added!\n")
      {:error, _changeset} -> send_resp(conn, 400, "Malformed request!\n")
    end
  end

  post "/addLock" do
    send_resp(conn, 200, "/addLock")
  end

  get "/listUsers" do
    send_resp(conn, 200, "/listUsers")
  end

  get "/listRoles" do
    send_resp(conn, 200, "/listRoles")
  end

  get "/listLocks" do
    send_resp(conn, 200, "/listLocks")
  end

  match _ do
    conn |> send_resp(404, "Bad endpoint.\n")
  end
end
