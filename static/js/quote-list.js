class QuoteItem extends HTMLElement {
  connectedCallback() {
    this.content = this.getAttribute("content");
    this.createdBy = this.getAttribute("created-by");
    this.createdAt = this.getAttribute("created-at");
    this.canEditQuote = this.getAttribute("can-edit-quote") === "true";
    this.quoteId = this.getAttribute("quote-id");

    console.log({
      content: this.content,
      createdBy: this.createdBy,
      createdAt: this.createdAt,
      canEditQuote: this.canEditQuote,
      quoteId: this.quoteId,
    })

    this.innerHTML = `
      <div class="card bg-base-100 mb-4">
        <div class="card-body">
          <h2 class="card-title mb-2">${this.content}</h2>
          <div class="flex justify-between items-end">
            <ul class="flex flex-col gap-2">
              <li class="flex gap-1 opacity-70 items-center">
                <custom-icon type="user"></custom-icon>
                <span>${this.createdBy}</span>
              </li>
              <li class="flex gap-1 opacity-70 items-center">
                <custom-icon type="date"></custom-icon>
                <span>${this.createdAt}</span>
              </li>
            </ul>
            ${this.canEditQuote ? `
              <div class="flex gap-2">
                <a
                  href="/quotes/${this.quoteId}/edit"
                  class="btn btn-square btn-outline btn-accent flex items-center justify-center"
                >
                                  <custom-icon type="edit"></custom-icon>
                </a>
                <button
                  hx-confirm="Are you sure you want to delete this quote?"
                  data-confirm-title="Delete quote"
                  hx-delete="/quotes/${this.quoteId}"
                  hx-target="closest .card"
                  hx-swap="outerHTML"
                  class="btn btn-square btn-outline btn-error flex items-center justify-center"
                >
                                  <custom-icon type="delete"></custom-icon>
                </button>
              </div>
              `
              :
              ""
            }
          </div>
        </div>
      </div>
    `;
  }
}

customElements.define('quote-item', QuoteItem);
