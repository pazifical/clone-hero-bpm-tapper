let tapTimes = [];
const audio = document.getElementById("audio");
const tapButton = document.getElementById("tap");
const playButton = document.getElementById("play");
const waveform = document.getElementById("waveform");
document.getElementById("play").focus();

function play() {
  audio.currentTime = 0;
  // tapTimes.push(0);
  audio.play();
  tapButton.focus();
}

async function stop() {
  audio.pause();

  let response = await fetch("/api/taps", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(tapTimes),
  });

  document.getElementById(
    "info"
  ).innerHTML = `First tap: <strong>${tapTimes[0]}</strong> seconds.`;
  console.log(response);
}

function reset() {
  stop();
  audio.currentTime = 0;
  tapTimes = [];
  renderTapTimes();
  playButton.focus();
}

function tap() {
  tapTimes.push(audio.currentTime);
  renderTapTimes();
}

function calcBPM() {
  calcAverageBPM();
  calcPartsBPM();
}

function calcPartsBPM() {
  console.log(tapTimes);
  const n = Number(document.getElementById("taps-average").value);
  const listItems = [];
  let i = 0;
  while (i < tapTimes.length - n) {
    const currentTapTimes = tapTimes.slice(i, i + n - 1);
    console.log(i, i + n - 1);

    const starts = currentTapTimes.slice(0, -1);
    const ends = currentTapTimes.slice(1);
    console.log(starts, ends);

    const dts = [];
    for (let j = 0; j < starts.length; j++) {
      dts.push(ends[j] - starts[j]);
    }
    const avg =
      dts.reduce((a, b) => {
        return a + b;
      }, 0) / dts.length;
    console.log("avg", avg);

    const t0 = tapTimes[i];
    const t1 = tapTimes[i + n - 1];

    const bpm = (1 / avg) * 60;
    const bpmRounded = Math.round(bpm * 100) / 100;
    listItems.push(`<li>${bpmRounded} BPM (${t0}s to ${t1}s)</li>`);
    i += n - 1;
  }
  document.getElementById("bpms").innerHTML = listItems.join("");
}

// async function download(event) {
//   const response = fetch("/api/charts", {
//     method: "POST",
//     headers: {
//       "Content-Type": "application/json",
//     },
//     body: JSON.stringify({
//       name: document.getElementById("name").value,
//       artist: document.getElementById("artist").value,
//     }),
//   });
//   event.preventDefault();
// }

document.getElementById("download-form").onsubmit = async (event) => {
  const form = event.target;
  const formData = new FormData(form);
  const response = await fetch("/api/charts", {
    method: "POST",
    body: formData,
  });
  event.preventDefault();
  return false;
};

function calcAverageBPM() {
  if (tapTimes.length < 2) {
    return;
  }
  const dt = tapTimes[tapTimes.length - 1] - tapTimes[0];
  const bpm = (tapTimes.length / dt) * 60;
  const bpmRounded = Math.round(bpm * 100) / 100;
  document.getElementById("bpm").value = bpmRounded;
}

function renderTapTimes() {
  waveform.innerHTML = tapTimes
    .map(
      (time) =>
        `<div class="marker" style="left: ${
          (time / audio.duration) * waveform.offsetWidth
        }px;"></div>`
    )
    .join("");
}

async function updateSongs() {
  const response = await fetch("/api/songs");
  const songList = await response.json();
  const songSelector = document.getElementById("songs");
  songSelector.innerHTML =
    "<option></option>" +
    songList.map((song) => `<option value="${song}">${song}</option>`).join("");
}

audio.ontimeupdate = () => {};
document.getElementById("songs").onchange = (e) => {
  const songName = e.target.value;
  console.log(songName);
  document.getElementById("audio").src = `/songs/${songName}`;
};

window.onload = async () => {
  await updateSongs();
};

document.getElementById("upload-file").onchange = async (event) => {
  const form = document.getElementById("upload-form");
  const formData = new FormData(document.getElementById("upload-form"));
  const response = await fetch("/api/songs", {
    method: "POST",
    body: formData,
  });
  if (response.ok) {
    document.getElementById(
      "upload-info"
    ).innerText = `Successfully uploaded ${event.target.value}`;
    event.target.value = "";
  }

  event.preventDefault();
};
