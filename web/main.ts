import { html, render } from "lit-html";

type ServerData = { keyvalues: any; log: [RawLogLine] };
type RawLogLine = { timestamp: number; query: string };

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

const Main = async () => {
  const data: ServerData = await (await fetch("/data/json")).json();
  console.log(data);

  render(MainElement(data), document.body);
};

Main();
