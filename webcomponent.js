class extends HTMLElement {
	constructor() {
		super();

		this.setRGBA = this.setRGBA.bind(this);
		this.sliderUpdate = this.sliderUpdate.bind(this);

		this.r = 120;
		this.g = 120;
		this.b = 120;
		this.a = 120;

		this.attachShadow({mode: "open"});
		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/standard.css"/>
<style>
#slider_holder {
	display: flex;
	height: 12em;
	overflow: hidden;
	position: relative;
}
#slider_holder > input {
	-webkit-appearance: none;
	transform: rotate(270deg);
	transform-origin: 50% 50%;
	width: 10em;
	height: 0.2em;
	top: 3em;
	position: absolute;
	margin-top: 3em;
	background-color: var(--clr-red);
}
</style>
<div class="jmod-wrapper">
	<div class="jmod-header" style="display:flex">
		<h1>RGB Light</h1>
		<svg viewBox="0 0 360 360">
			<path d="M150,300 A30,60,0,0,1,210,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M105,300 A60,90,0,0,1,255,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M60,300 A60,80,0,0,1,300,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
		</svg>
	</div>
	<hr>
	<div class="jmod-body">
		<div id="debug-out">
		</div>
		<p>This will have sliders for controlling the color and a selector for which lights to send the set command to</p>

		<div id="slider_holder">
			<input type="range" id="slider_r" hmin="0" max="255" value="120" oninput="this.getRootNode().host.sliderUpdate('r', parseInt(this.value, 10))" style="left: -10%;"></input>
			<input type="range" id="slider_g" hmin="0" max="255" value="120" oninput="this.getRootNode().host.sliderUpdate('g', parseInt(this.value, 10))" style="left: 15%; background-color: var(--clr-green);"></input>
			<input type="range" id="slider_b" hmin="0" max="255" value="120" oninput="this.getRootNode().host.sliderUpdate('b', parseInt(this.value, 10))" style="left: 40%; background-color: var(--clr-blue);"></input>
			<input type="range" id="slider_a" hmin="0" max="255" value="120" oninput="this.getRootNode().host.sliderUpdate('a', parseInt(this.value, 10))" style="left: 70%; background-color: var(--clr-font-med);"></input>
		</div>
	</div>
</div>
		`
	}

	init(source, instName, config) {
		this.source = source;
		this.instName = instName;
		this.config = config;
		this.socket = new WebSocket(`ws://${document.location.host}/jmod/socket?JMOD-Source=${this.source}`);

		const elem = this.shadowRoot.getElementById("debug-out");
		for (var x of this.config.lightIPs) {
			elem.textContent += "IP:" + x + "\n";
		}
	}

	sliderUpdate(color, value) {
		this[color]	= value
		this.setRGBA();
	}

	setRGBA() {
		this.socket.send(`${this.config.lightIPs[1]},${this.r},${this.g},${this.b},${this.a}`);
		console.log(this.r, this.g, this.b, this.a);
	}
}
