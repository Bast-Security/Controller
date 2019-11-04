defmodule Bast.Middleman.REST do
  @moduledoc """
  This is the HTTP server responsible for handling HTTP
  REST messages sent to the middleman.
  """
  
  alias Bast.Middleman.{User}

  use Plug.Router
  plug Plug.Logger
  plug Plug.Parsers,
    parsers: [:json],
    json_decoder: Jason
  plug :formatter
  plug :match
  plug :dispatch

  post "/addUser" do
    user = struct(User, conn.body_params)
    {:ok, _user} = Bast.Middleman.Repo.insert user
    conn |> send_resp(200, "User added!\n")
  end

  match _ do
    conn |> send_resp(404, "Bad endpoint.\n")
  end

  defp formatter(conn, _opts) do
    new_params = Map.new(conn.body_params, fn {k, v} -> {String.to_atom(k), v} end)
    %{conn | body_params: new_params}
  end
end
