defmodule Bast.Api.Admins do
  use Plug.Router
  alias Bast.Repo
  alias Bast.Api.Status
  alias Bast.Admin

  plug :match
  plug :dispatch

  post "/register" do
    inserted =
      conn.body_params
      |> Admin.create
      |> Repo.insert

    status =
      case inserted do
        {:ok, _struct} -> Status.created
        {:error, _changeset} -> Status.malformed_request
      end

    Status.send_status(conn, status)
  end

  match _ do
    Status.send_status(conn, Status.not_found)
  end
end
