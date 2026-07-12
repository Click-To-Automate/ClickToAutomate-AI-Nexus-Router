import { useEffect, useRef, useState } from "react";
import gsap from "gsap";
import './Preloader.css';

const TAGLINE = "Websites · Marketing · Automation";
const WORD = "Automate".split("");

// Node Icons
const LeadCaptureIcon = () => (
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" width={16} height={16}><path d="M5 12a7 7 0 0 1 7-7v2a5 5 0 0 0-5 5H5z"/><path d="M12 5a7 7 0 0 1 7 7h-2a5 5 0 0 0-5-5V5z"/><circle cx="12" cy="12" r="2"/></svg>
);
const CrmIcon = () => (
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" width={16} height={16}><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/><path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3"/></svg>
);
const AiIcon = () => (
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" width={16} height={16}><rect x="4" y="4" width="16" height="16" rx="2"/><path d="M9 9h6v6H9z"/><path d="M9 1v3M15 1v3M9 20v3M15 20v3M20 9h3M20 15h3M1 9h3M1 15h3"/></svg>
);
const EmailIcon = () => (
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" width={16} height={16}><rect x="2" y="4" width="20" height="16" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
);
const PipelineIcon = () => (
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" width={16} height={16}><polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/><polyline points="16 7 22 7 22 13"/></svg>
);
const GrowthIcon = () => (
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" width={16} height={16}><path d="M4.5 16.5c-1.5 1.26-2 2.6-1.5 3 1 .8 2.4.5 3.5-.5L9 16.5H4.5z"/><path d="M12 15 9 16.5 6 13.5 7.5 10.5 12 15z"/><path d="M11 14.5c.3-1.6 1.3-4.5 4.5-7.5s5.9-4.2 7.5-4.5c.3 1.6-.5 4.3-3.7 7.5s-5.9 4.2-7.5 4.5z"/><path d="M13.5 9a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3z"/></svg>
);

const NODES = [
  { id: 1, label: "Lead Capture", x: 120, y: 100, Icon: LeadCaptureIcon, align: "left" },
  { id: 2, label: "CRM", x: 280, y: 220, Icon: CrmIcon, align: "right" },
  { id: 3, label: "AI Processing", x: 120, y: 340, Icon: AiIcon, align: "left" },
  { id: 4, label: "Email Automation", x: 280, y: 460, Icon: EmailIcon, align: "right" },
  { id: 5, label: "Sales Pipeline", x: 120, y: 580, Icon: PipelineIcon, align: "left" },
  { id: 6, label: "Business Growth", x: 200, y: 700, Icon: GrowthIcon, align: "center" }
];

const PATHS = [
  { id: 1, d: "M 120 100 C 120 160, 280 160, 280 220" },
  { id: 2, d: "M 280 220 C 280 280, 120 280, 120 340" },
  { id: 3, d: "M 120 340 C 120 400, 280 400, 280 460" },
  { id: 4, d: "M 280 460 C 280 520, 120 520, 120 580" },
  { id: 5, d: "M 120 580 C 120 640, 200 640, 200 700" }
];

const ALL_PATH_D = "M 120 100 C 120 160, 280 160, 280 220 C 280 280, 120 280, 120 340 C 120 400, 280 400, 280 460 C 280 520, 120 520, 120 580 C 120 640, 200 640, 200 700";

interface PreloaderProps {
  onComplete?: () => void;
}

export function Preloader({ onComplete }: PreloaderProps) {
  const [complete, setComplete] = useState(false);
  const [shouldRender, setShouldRender] = useState(false);
  const [particleCount, setParticleCount] = useState(15);
  
  const containerRef = useRef<HTMLDivElement>(null);
  const curtainRef = useRef<HTMLDivElement>(null);
  const svgGroupRef = useRef<SVGGElement>(null);
  const pulseRef = useRef<SVGPathElement>(null);
  
  const nodeRefs = useRef<(SVGGElement | null)[]>([]);
  const nodeGlowRefs = useRef<(SVGCircleElement | null)[]>([]);
  const labelRefs = useRef<(SVGTextElement | null)[]>([]);
  const pathRefs = useRef<(SVGPathElement | null)[]>([]);
  
  const textWordsContainerRef = useRef<HTMLDivElement>(null);
  const textWordRefs = useRef<(HTMLSpanElement | null)[]>([]);
  const particlesContainerRef = useRef<HTMLDivElement>(null);

  const clickToRef = useRef<HTMLHeadingElement>(null);
  const letterRefs = useRef<(HTMLSpanElement | null)[]>([]);
  const letterGlowRefs = useRef<(HTMLSpanElement | null)[]>([]);
  const taglineRef = useRef<HTMLParagraphElement>(null);
  const introTaglineRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    const hasSeen = sessionStorage.getItem("splashSeen");
    if (hasSeen === "true") {
      setComplete(true);
      if (onComplete) onComplete();
    } else {
      setShouldRender(true);
      setParticleCount(window.innerWidth < 640 ? 8 : 18);
    }
  }, [onComplete]);

  useEffect(() => {
    if (!shouldRender) return;

    pathRefs.current.forEach((path) => {
      if (path) {
        const len = path.getTotalLength();
        gsap.set(path, { strokeDasharray: len, strokeDashoffset: len });
      }
    });

    if (pulseRef.current) {
      const pulseLen = pulseRef.current.getTotalLength();
      gsap.set(pulseRef.current, { strokeDasharray: `40 ${pulseLen}`, strokeDashoffset: pulseLen });
    }

    gsap.set(containerRef.current, { autoAlpha: 1 });
    gsap.set(nodeRefs.current, { scale: 0, opacity: 0, transformOrigin: "center" });
    gsap.set(nodeGlowRefs.current, { opacity: 0 });
    gsap.set(labelRefs.current, { fill: "rgba(0, 0, 0, 0.35)" });
    gsap.set(textWordRefs.current, { opacity: 0, y: 15, filter: "blur(8px)" });

    gsap.set(clickToRef.current, { opacity: 0, x: -45, filter: "blur(10px)" });
    gsap.set(taglineRef.current, { opacity: 0, y: 12, letterSpacing: "0.1em" });
    gsap.set(letterRefs.current, { opacity: 0, scale: 0, rotate: () => gsap.utils.random(-80, 80) });
    gsap.set(letterGlowRefs.current, { opacity: 0, scale: 0 });
    gsap.set(introTaglineRef.current, { opacity: 0, y: 15, filter: "blur(8px)" });

    gsap.set(svgGroupRef.current, { opacity: 1 });

    const tl = gsap.timeline({
      onComplete: () => {
        sessionStorage.setItem("splashSeen", "true");
        setComplete(true);
        if (onComplete) onComplete();
      },
    });

    const particles = particlesContainerRef.current?.children;
    if (particles) {
      Array.from(particles).forEach((p) => {
        gsap.to(p, {
          y: `-=${gsap.utils.random(40, 100)}`,
          x: `+=${gsap.utils.random(-30, 30)}`,
          opacity: gsap.utils.random(0.05, 0.2),
          duration: gsap.utils.random(8, 15),
          repeat: -1,
          yoyo: true,
          ease: "sine.inOut"
        });
      });
    }

    NODES.forEach((_, i) => {
      const nodeTime = i * 0.22;
      tl.to(nodeRefs.current[i], {
        scale: 1,
        opacity: 1,
        duration: 0.65,
        ease: "elastic.out(1, 0.75)"
      }, nodeTime);

      if (i < PATHS.length) {
        tl.to(pathRefs.current[i], {
          strokeDashoffset: 0,
          duration: 0.65,
          ease: "power3.inOut"
        }, nodeTime + 0.12);
      }
    });

    tl.to(pulseRef.current, {
      strokeDashoffset: 0,
      duration: 1.3,
      ease: "power2.inOut"
    }, 1.3);

    NODES.forEach((_, i) => {
      const triggerTime = 1.3 + (i * 0.2);
      tl.to(nodeGlowRefs.current[i], { opacity: 1, duration: 0.3, ease: "power2.out" }, triggerTime);
      tl.to(labelRefs.current[i], { fill: "#000000", duration: 0.3, ease: "power2.out" }, triggerTime);
      tl.to(nodeRefs.current[i], { scale: 1.08, duration: 0.18, yoyo: true, repeat: 1 }, triggerTime);
    });

    tl.addLabel("collapse", 2.6);
    
    NODES.forEach((node, i) => {
      tl.to(nodeRefs.current[i], {
        x: 200 - node.x,
        y: 390 - node.y,
        scale: 0.08,
        opacity: 0,
        duration: 0.75,
        ease: "power4.inOut"
      }, "collapse");
    });

    tl.to(svgGroupRef.current, {
      scale: 0.12,
      opacity: 0,
      transformOrigin: "200px 390px",
      duration: 0.75,
      ease: "power4.inOut"
    }, "collapse");

    tl.to(textWordRefs.current[0], { opacity: 1, y: 0, filter: "blur(0px)", duration: 0.45, ease: "power3.out" }, 3.0);
    tl.to(textWordRefs.current[1], { opacity: 1, y: 0, filter: "blur(0px)", duration: 0.45, ease: "power3.out" }, 3.2);
    tl.to(textWordRefs.current[2], { opacity: 1, y: 0, filter: "blur(0px)", duration: 0.45, ease: "power3.out" }, 3.4);

    tl.to(textWordsContainerRef.current, { opacity: 0, scale: 0.95, filter: "blur(8px)", duration: 0.4, ease: "power3.in" }, 4.0);

    tl.to(introTaglineRef.current, { opacity: 1, y: 0, filter: "blur(0px)", duration: 0.65, ease: "power3.out" }, 4.25);
    tl.to(introTaglineRef.current, { opacity: 0, y: -15, filter: "blur(8px)", duration: 0.5, ease: "power3.in" }, 5.25);

    tl.to(clickToRef.current, { opacity: 1, x: 0, filter: "blur(0px)", duration: 0.85, ease: "power4.out" }, 5.65);

    tl.to(letterRefs.current, {
      opacity: 1, scale: 1, rotate: 0, duration: 0.85, stagger: 0.03, ease: "elastic.out(1.2, 0.6)"
    }, 5.7);

    tl.to(letterGlowRefs.current, {
      opacity: 0.6, scale: 1.1, duration: 0.85, stagger: 0.03, ease: "elastic.out(1.2, 0.6)"
    }, 5.7);

    tl.to(taglineRef.current, {
      opacity: 1, y: 0, letterSpacing: "0.25em", duration: 0.75, ease: "power3.out"
    }, 5.9);

    tl.to([clickToRef.current, ...letterRefs.current, ...letterGlowRefs.current, taglineRef.current], {
      y: -40, opacity: 0, scale: 0.98, filter: "blur(5px)", duration: 0.55, stagger: 0.01, ease: "power4.in",
    }, 7.2)
    .to(curtainRef.current, { yPercent: -100, duration: 0.8, ease: "expo.inOut" }, 7.4)
    .to(containerRef.current, { autoAlpha: 0, duration: 0.1 });

    return () => { tl.kill(); };
  }, [shouldRender]);

  if (complete) return null;

  return (
    <div ref={containerRef} className="preloader-overlay">
      <div ref={curtainRef} className="preloader-curtain">
        {shouldRender && (
          <div ref={particlesContainerRef} className="preloader-particles">
            {Array.from({ length: particleCount }).map((_, i) => (
              <div
                key={i}
                className="preloader-particle"
                style={{
                  width: `${Math.random() * 3 + 1}px`,
                  height: `${Math.random() * 3 + 1}px`,
                  left: `${Math.random() * 100}%`,
                  top: `${Math.random() * 100}%`,
                  opacity: Math.random() * 0.15 + 0.05
                }}
              />
            ))}
          </div>
        )}

        <div className="preloader-grid" />
        <div className="preloader-glow" />

        {shouldRender && (
          <div className="preloader-svg-container">
            <svg className="preloader-svg" viewBox="0 0 400 800">
              <defs>
                <filter id="shadow" x="-30%" y="-30%" width="160%" height="160%">
                  <feDropShadow dx="0" dy="4" stdDeviation="5" floodColor="#ff6a3d" floodOpacity="0.2" />
                </filter>
                <linearGradient id="path-glow-gradient" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stopColor="#ff6a3d" />
                  <stop offset="50%" stopColor="#27f1f7" />
                  <stop offset="100%" stopColor="#87af32" />
                </linearGradient>
                <linearGradient id="pulse-gradient" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stopColor="#ff8c6b" />
                  <stop offset="50%" stopColor="#27f1f7" />
                  <stop offset="100%" stopColor="#87af32" />
                </linearGradient>
              </defs>

              <g ref={svgGroupRef} style={{ opacity: 0, transformOrigin: "200px 390px" }}>
                {PATHS.map((path) => (
                  <path key={`base-path-${path.id}`} d={path.d} stroke="rgba(0, 0, 0, 0.04)" strokeWidth="1.5" fill="none" />
                ))}

                {PATHS.map((path, idx) => (
                  <path key={`active-path-${path.id}`} ref={(el) => { pathRefs.current[idx] = el; }} d={path.d} stroke="url(#path-glow-gradient)" strokeWidth="2" fill="none" />
                ))}

                <path ref={pulseRef} d={ALL_PATH_D} stroke="url(#pulse-gradient)" strokeWidth="3.5" fill="none" strokeLinecap="round" filter="url(#shadow)" />

                {NODES.map((n, i) => {
                  const IconComponent = n.Icon;
                  return (
                    <g key={`node-${n.id}`} ref={(el) => { nodeRefs.current[i] = el; }} style={{ transformOrigin: `${n.x}px ${n.y}px` }}>
                      <circle cx={n.x} cy={n.y} r="24" fill="#ffffff" stroke="rgba(0, 0, 0, 0.06)" strokeWidth="1.5" style={{ filter: "drop-shadow(0 4px 10px rgba(0,0,0,0.03))" }} />
                      <circle ref={(el) => { nodeGlowRefs.current[i] = el; }} cx={n.x} cy={n.y} r="24" fill="transparent" stroke="url(#pulse-gradient)" strokeWidth="2" filter="url(#shadow)" />
                      <g transform={`translate(${n.x - 8}, ${n.y - 8})`} stroke="#111827">
                        <IconComponent />
                      </g>
                      <text ref={(el) => { labelRefs.current[i] = el; }} x={n.align === "left" ? n.x + 36 : n.align === "right" ? n.x - 36 : n.x} y={n.align === "center" ? n.y + 42 : n.y + 5} textAnchor={n.align === "left" ? "start" : n.align === "right" ? "end" : "middle"} className="preloader-label">
                        {n.label}
                      </text>
                    </g>
                  );
                })}
              </g>
            </svg>
          </div>
        )}

        <div ref={textWordsContainerRef} className="preloader-text-container">
          <div className="preloader-text-inner">
            <span ref={(el) => { textWordRefs.current[0] = el; }} className="preloader-word-1">Automate.</span>
            <span ref={(el) => { textWordRefs.current[1] = el; }} className="preloader-word-2">Scale.</span>
            <span ref={(el) => { textWordRefs.current[2] = el; }} className="preloader-word-3">Grow.</span>
          </div>
        </div>

        <div className="preloader-intro-container">
          <h2 ref={introTaglineRef} className="preloader-intro-tagline">
            Run your business on autopilot with
          </h2>
        </div>

        <div className="preloader-logo-container">
          <div className="preloader-logo-inner">
            <div className="preloader-logo-flex">
              <h1 ref={clickToRef} className="preloader-logo-prefix">ClickTo</h1>
              <div className="preloader-logo-suffix">
                {WORD.map((char, index) => (
                  <div key={index} className="preloader-char-wrapper">
                    <span ref={el => { letterGlowRefs.current[index] = el; }} className="preloader-char-glow">{char}</span>
                    <span ref={el => { letterRefs.current[index] = el; }} className="preloader-char-solid">{char}</span>
                  </div>
                ))}
              </div>
            </div>
            <p ref={taglineRef} className="preloader-logo-tagline">{TAGLINE}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
