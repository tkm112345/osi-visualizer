import type { Step } from "../types";
import LayerCard from "./LayerCard";

interface Props {
  steps: Step[];
  selectedLevel: number | null;
  activeLevel: number | null;
  onSelect: (level: number) => void;
}

// steps は L7 → L1 の順で渡ってくる。積み上げ表示なので上から L7 で問題なし。
export default function LayerStack({ steps, selectedLevel, activeLevel, onSelect }: Props) {
  // アニメーションは L7 から L1 へ進む。activeLevel 以上（level が大きい方）は到達済み。
  const reachedThreshold = activeLevel ?? 1;

  return (
    <div className="layer-stack">
      {steps.map((step) => (
        <LayerCard
          key={step.level}
          step={step}
          selected={selectedLevel === step.level}
          active={activeLevel === step.level}
          reached={activeLevel === null || step.level >= reachedThreshold}
          onClick={() => onSelect(step.level)}
        />
      ))}
    </div>
  );
}
