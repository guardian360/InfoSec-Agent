class Navbar extends HTMLElement {
    constructor() {
      super();
    }
  
    connectedCallback() {
      this.innerHTML = `
        <header>
          <nav>
            <ul>
              <li><a href="about.html">About</a></li>
              <li><a href="work.html">Work</a></li>
              <li><a href="contact.html">Contact</a></li>
            </ul>
          </nav>
        </header>
      `;
    }
  }
  
  customElements.define('vertical-navbar', Navbar);

