import { html, render } from "lit-html";

const UPDATE_PERIOD = 100;

type ServerData = { keyvalues: any; log: [RawLogLine] };
type RawLogLine = { timestamp: number; query: string };

let data: ServerData;

const QUERY_COLORS = {
  LOG: "text-blue-500",
  _: "text-red-600",
};

const ParseTimeStamp = (timestamp: number) => {
  const hh = Math.floor(timestamp / 3600_000)
    .toString()
    .padStart(2, "0");
  const mm = (Math.floor(timestamp / 60_000) % 60).toString().padStart(2, "0");
  const ss = (Math.floor(timestamp / 1000) % 60).toString().padStart(2, "0");
  const mss = (timestamp % 1000).toString().padStart(3, "0");
  return `${hh}:${mm}:${ss}.${mss}`;
};

const QuerySplit = (query: string) => {
  const [cmd, args] = query.split(" ", 2);
  const color = QUERY_COLORS[cmd] || QUERY_COLORS["_"];
  return {
    cmd,
    query: html`<span class="${color}">${cmd}</span> <span>${args}</span>`,
  };
};

const MainElement = ({ keyvalues, log }: ServerData) =>
  html`<div class="font-mono">
    <table>
      ${log.map(LogElement)}
    </table>
  </div>`;

const LogElement = (log: RawLogLine) =>
  html`<tr>
    <td class="p-0">[${ParseTimeStamp(log.timestamp)}]</td>
    <td class="p-0">${QuerySplit(log.query).query}</td>
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
