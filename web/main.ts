import { html, render } from "lit-html";

const UPDATE_PERIOD = 100;

type ServerData = { keyvalues: { [key: string]: string }; log: [RawLogLine] };

type RawLogLine = {
  timestamp: number;
  query: string;
  attributes: { group: string };
};

let data: ServerData;

const CMD_COLORS = {
  LOG: "text-blue-500",
  _: "text-red-600",
};

const ARGS_COLORS = {
  invalid_cmd: "text-red-600",
  _: "text-black",
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

const MainElement = ({ keyvalues, log }: ServerData) =>
  html`<div class="font-mono">
    <table>
      ${log.map(LogElement)}
    </table>
  </div>`;

const SplitQuery = (query: string): [string, string] => {
  let i = 0;
  while (i < query.length && query[i] != " ") {
    i += 1;
  }
  return [query.slice(0, i), query.slice(i + 1)];
};

const LogElement = (log: RawLogLine) => {
  const [cmd, args] = SplitQuery(log.query);
  const cmd_color = CMD_COLORS[cmd] || CMD_COLORS["_"];
  const args_color = ARGS_COLORS[log.attributes.group] || ARGS_COLORS["_"];

  return html`<tr>
    <td>[${ParseTimeStamp(log.timestamp)}]</td>
    <td>
      <span class="${cmd_color}">${cmd}</span>
      <span class="${args_color}">${args}</span>
    </td>
  </tr>`;
};

const Render = async () => {
  data = await (await fetch("/data/json")).json();
  render(MainElement(data), document.body);
  setTimeout(Render, UPDATE_PERIOD);
};

const Main = async () => {
  Render();
};

Main();
