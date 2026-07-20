import type { Step } from "../types";
import LayerCard from "./LayerCard";

interface Props {
  steps: Step[];
  mode: "encap" | "decap";
  selectedLevel: number | null;
  activeLevel: number | null;
  onSelect: (level: number) => void;
}

// 常に L7（上）〜 L1（下）で表示する。
// encap: データは上→下へ進むので activeLevel 以上（level が大きい方）が到達済み。
// decap: データは下→上へ進むので activeLevel 以下（level が小さい方）が到達済み。
export default function LayerStack({ steps, mode, selectedLevel, activeLevel, onSelect }: Props) {
  const display = [...steps].sort((a, b) => b.level - a.level); // L7 → L1

  function isReached(level: number): boolean {
    if (activeLevel === null) return true;
    return mode === "encap" ? level >= activeLevel : level <= activeLevel;
  }

  return (
    <div className="layer-stack">
      {display.map((step) => (
        <LayerCard
          key={step.level}
          step={step}
          selected={selectedLevel === step.level}
          active={activeLevel === step.level}
          reached={isReached(step.level)}
          onClick={() => onSelect(step.level)}
        />
      ))}
    </div>
  );
}
