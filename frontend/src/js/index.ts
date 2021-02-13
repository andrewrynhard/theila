import "./style.css";

// Turbo
const Turbo = require("@hotwired/turbo");
Turbo.connectStreamSource(
  new WebSocket("ws://" + document.location.host + "/servers/socket")
);
Turbo.connectStreamSource(
  new WebSocket("ws://" + document.location.host + "/clusters/socket")
);

// Stimulus
import { Application } from "stimulus";
import { definitionsFromContext } from "stimulus/webpack-helpers";

const application = Application.start();
const context = require.context("./controllers", true, /\.ts$/);
application.load(definitionsFromContext(context));
