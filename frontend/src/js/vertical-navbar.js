class Navbar extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.innerHTML = `
    <div id="header">
      <div id="header-logo">
          <a href="/home.html" class="logo-name">
              <img src="./src/assets/images/logoTeamA-transformed.png" alt="logo"></img>
              <div id="header-name">
                  <h1>Little Brother</h1>
              </div>
          </a>
      </div>
      <div id="header-settings">
          <a href="./settings.html"><p><span class="material-symbols-outlined">settings</span><span>Settings</span></p></a>
      </div>
    </div> 
    <div id="left-nav">
        <a href="./index.html">
            <p><span class="material-symbols-outlined">home</span><span class="nav-item">Home</span></p>
        </a>
        <a href="./dashboard.html">
            <p><span class="material-symbols-outlined">monitoring</span><span class="nav-item">Dashboard</span></p>
        </a>
        <a href="./issues.html">
            <p><span class="material-symbols-outlined">security</span><span class="nav-item">Issues</span></p>
        </a>
        <a href="./integration.html">
            <p><span class="material-symbols-outlined">integration_instructions</span><span class="nav-item">Integration</span></p>
        </a>
        <a href="./about.html">
            <p><span class="material-symbols-outlined">info</span><span class="nav-item">About</span></p>
        </a>
    </div>
    `;
  }
}

customElements.define('vertical-navbar', Navbar);