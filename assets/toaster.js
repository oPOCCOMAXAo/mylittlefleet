function showToasts() {
  document.querySelectorAll(".toast").forEach((t) => new bootstrap.Toast(t).show());
}
