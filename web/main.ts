import { html, render } from "lit-html";

const UPDATE_PERIOD = 100;

type ServerData = { keyvalues: any; log: [RawLogLine] };
type RawLogLine = { timestamp: number; query: string };

let data: ServerData;

const MainElement = ({ keyvalues, log }: ServerData) =>
  html`<div class="font-mono">
    <table>
      ${log.map(LogElement)}
    </table>
  </div>`;

const LogElement = (log: RawLogLine) =>
  html`<tr>
    <td>[${log.timestamp}]</td>
    <td>${log.query}</td>
  </tr>`;

const Render = async () => {
  data = await (await fetch("/data/json")).json();
  render(MainElement(data), document.body);
  setTimeout(Render, UPDATE_PERIOD);
};

const Main = async () => {
  Render();
};

Main();
