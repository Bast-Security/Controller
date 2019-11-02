defmodule Bast.Middleman.REST do
  @moduledoc """
  This is the HTTP server responsible for handling HTTP
  REST messages sent to the middleman.
  """

  use Plug.Router

  plug :match
  plug :dispatch

  get "/" do
    conn |> send_resp(200, "Hello World!\n")
  end

  match _ do
    conn |> send_resp(404, "Bad endpoint.\n")
  end
end
