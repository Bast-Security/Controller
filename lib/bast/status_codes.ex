defmodule Bast.Api.Status do
  def created, do: { 201, "Created\n" }

  def malformed_request, do: { 422, "Malformed Request.\n" }

  def unauthorized, do: { 401, "Unauthorized\n" }

  def forbidden, do: { 403, "Forbidden\n" }

  def not_found, do: { 404, "Not Found\n" }

  def not_implemented, do: { 501, "Not Implemented\n" }

  def send_status(conn, {status, body}) do
    Plug.Conn.send_resp(conn, status, body)
  end
end
