defmodule Bast.Controller.REST do
  @moduledoc """
  This is the HTTP server responsible for handling HTTP
  REST messages sent to the middleman.
  """

  alias Bast.Controller.User

  use Plug.Router

  plug Plug.Logger
  plug Plug.Parsers,
    parsers: [:json],
    json_decoder: Jason
  plug :formatter
  plug :match
  plug :dispatch

  post "/addUser" do
    user = struct(User, conn.body_params) |> User.validate
    case Bast.Controller.Repo.insert(user, on_conflict: :nothing) do
      {:ok, _user} -> send_resp(conn, 200, "User added!\n")
      {:error, _changeset} -> send_resp(conn, 400, "Malformed request!\n")
    end
  end

  post "/addRole" do
    send_resp(conn, 200, "/addRole")
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

  defp formatter(conn, _opts) do
    new_params = Map.new(conn.body_params, fn {k, v} -> {String.to_atom(k), v} end)
    %{conn | body_params: new_params}
  end
end
