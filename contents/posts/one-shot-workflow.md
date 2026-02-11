---
title: Context Rules Everything Around Me
date: 2026-02-04
last-modified: 2026-02-11
subtitle: I'm tired of screwing around in a TUI.
---

I am a little bit disappointed with my industry right now. There seems to be have been a massive unmooring of steady thought over the past year or so, brought on by the full advent of LLMs and highly competent coding agents. Much of this is driven by emotions from purported AI driven layoffs and belt tightening[^0], a mourning for the zen and art of coding, and a general identity crisis around who a software engineer is and what they do. Much of it is also a lack of imagination and a devolution to the lowest common denominator now that there is no longer a clear status signal between those who carefully craft and those who shovel slop. I am sure to many who previously considered themselves fine artisans it is easy to give up and start flinging vibecoded LOC over the wall. But I cannot shake the fact that this wave just _does not feel like engineering_, and I refuse to believe that the engineering of software systems is dead.

Our industry has (notoriously) long struggled with the concepts of rigor, systems discipline, and formal standards... often swinging violently between overdoing it and completely ignoring anything that might resemble exactness. Taking a step back one can easily frame the recent LLM driven chaos as this finally coming home to root and many reaping the seeds they sowed long ago. The obsession with things like frameworks, leetcode, and arcane bit twiddling[^1] and disregard for basic, repeatable, and structured _engineering processes_ has set many up for a stunning disappointment. No longer does it matter how well you can wrangle that array when querying from MySQL. What matters is how you structure your work, how quickly you can develop a pipeline of repeatable processes, and how _reliable and high quality_ you can make the extruded bits.

With the benefit of hindsight, this is really nothing new. I know many have believed that the reason they earned their keep was because they created beautiful experiences: SaaS that was a joy to use, and APIs that inspired delight. But truly, what we as an industry largely do is allow companies to turn CapEx into OpEx while presenting a throat that they may strangle when things go wrong. This is doubly, if not [triply true today](https://www.theinformation.com/newsletters/applied-ai/anthropic-jpmorgan-seem-agree-ai-eating-enterprise-apps-yet). Features are commodity, and to a certain extent they might have always been so. What now matters most is how well you can operate your service. Indeed I am increasingly of the opinion that as the next couple of years progress customers will begin to gravitate towards companies and services that are highly reliable and consistent to use, both due to what I forsee as instability in current offerings from poor LLM usage and because feature parity will be so easy to accomplish as the cost of LOC continues to plummet. In many ways the advent of LLMs might actually bring about the programmer nirvana: where those with the most beautiful, maintainable code win.

So why, then, are we seeing such a marked decline in code quality from users of LLMs, and how can we fix this? How can we provide highly reliable services when we are not writing the code that runs them?

## Context Windows

At risk of outing myself as ignorant, I simply refuse to believe that people are submitting 25+ PRs a day[^2] of meaningful, maintainable, quality work. It is entirely possible that some might have codebases riddled with bugs, submitted by others, that they can then feed in an automated fashion to coding agents to generate this volume but that then begs more questions than it answers. I also refuse to believe that people are creating maintainable pieces of software by letting swarms of agents code for hours on end, generating tens of thousands of lines of code, and then reasonably grokking that body of work. True understanding is borne in repitition and action, the ability to read volumunous amounts of code and truly internalize that into a mental model is an ability that is extraordinarily rare and not at all sustainable or scalable. There is the argument that models will continue to get better, that we can accrue tech debt and rely on future iterations to bail us out. I doubt it. That does not mean I doubt models will get better, I simply doubt that context windows will get meaningfully large enough in the near future[^3]. Self-attention scales quadratically with sequence length because every token attends to every other token. This makes GPU memory a _hard constraint_ on context windows. Indeed, despite the heralding of million token length windows my personal experience has been that regardless of model, context cliffs are a very real thing and have remained consistently within the first couple hundred thousand tokens. Even if this was not true, LLM context windows are still tiny compared to the human equivalent... we'd need to see an exponential increase in functional size to compete, and even then humans will have the advantage of _understanding more deeply as context grows_, as opposed to transformer models behaving in a unimodal fashion[^4]. I firmly believe that to build high quality, maintainable code bases with LLMs tasks must be broken up into relatively small, very precise, highly structured pieces of work and fed into empty or nearly empty contexts.

Once you truly internalize how precious context is it becomes more and more apparent how important structuring and detailing out your prompt is. Much has been said about prompt engineering and I will not waste words there, but I do know some simple guidelines. There is no reason for you to be wasting context with things like tools to search files, to read code that doesn't apply, and to discover codebase structure. Any relevant files, relevant code locations, and relevant code structure should be explicitly detailed in the prompt. Even more advanced sub-agent usage with context compaction does not compare to a careful human enumeration of what is relevant to the task at hand. Unless you are building something entirely greenfield in a domain you know nothing about (and this should not be happening if you are doing this as your day job. And if you are, at the very least use the LLM to research things yourself and distill your knowledge) there is no reason for the LLM to be doing research or web searches when given a coding task. Any relevant information, best practices, or specific domain knowledge should either be in the prompt or moved to re-usable skills.

The better you fully understand the task you are asking the LLM to do before hand, and the better you explain and structure this task, the higher quality result you will get. This is due to both context and priming. In my experience the best prompts resemble very small one pager design docs, are ~250-400 words, and explicitly detail out data structures, data flow, and general code structures. In other words, the human still does the thinking, but treats the LLM as a highly competent code gen. I do not think this is absolutely neccessary in all circumstances, but I do believe this style of prompting generates the highest quality code and has several knock on benefits.

## Code Review

The first of these benefits is dramatically faster and higher quality code review.

## Context as Kingmaker

It is certainly very telling that within months of being given the power of infinite junior engineers[^0] many have redesigned themselves as Mayors of their own little Jira Board Fiefdom. This should come as no surprise to anyone who has seen a competent engineer promoted to manager in a disfunctional Big Tech organization. A tale as old as time, they immediately become the authority they previously resented: micromanaging project managers. Of course, this is largely due to the absolute wasteland of true leadership that is most Big Tech companies and complete abdication of manager mentorship responsibilities by most VPs and Directors[^1] but it does speak to a very common undercurrent of the SWE psyche.

Despite what most would tell you, software engineers have incredible authority at most software companies for a very simple reason: they hold the system in their head. Code, it's subsequent compilers, and the infrastructure and silicon it runs on are worth nothing in times of operational crisis. Likewise, that beautiful pure function that operates on a higher plane of abstraction than most dharmic yogis means very little when not put to functional use. Engineers are not valuable for their ability to type at 50 wpm on weird and loud keyboards, they are valuable because they understand how to fit business requirements into a nebulous _something_ that generates money. And yet so many do not understand this. They see themselves as slaves to the eternal drip of backlogs so you see many promoted to a level of incompetence. Being recognized for excellence at moving tickets across the board will train people to see that as their value, and increasing their value just means the ticket moving becomes ever so more important. Unfortunately this is the environment under which we have just promoted the vast majority of our software producing workforce.

Truly senior[^2] engineers understand that when they are put in positions of leadership the system in their head and how it fits into the bigger picture technically, organizationally, and monetarily is worth it's weight in gold. Their job at that point is to not only democratize that knowledge but foster ownership in others so that they may understand it at a deeper level. As always, your very first point of business when promoted is to backfill your old position. If we are now all senior, then why have we become so obsessed with hovering over the shoulders of a flat file full of integers? Iteratively coaxing the LLM to spit out the right incantation to meet the almighty Definition Of Done is a fast way to solve no problems, very slowly, and at great cost. Most of your Linear projects could be filed as evidence at The Hague, but you believe shotgunning those half assed tickets at codex-5.2-xhigh will actually result in something useful? This is nothing but a distinguished sort of LLM psychosis[^3].

It drives me absolutely nuts to see how little thought and effort most engineers put into their prompts now. This I blame on the junk food nature of "vibecoding", the surprising intelligence of LLMs, and most of all a lack of understanding of how context windows work. The [topic](https://x.com/karpathy/status/1937902205765607626) of context and [it's importance](https://simonwillison.net/2025/Mar/11/using-llms-for-code/#context-is-king) is at this point almost overwrought[^4], and I know I sound like a broken record to those around me, but it is _extremely_ worth your while to [write your own agent](https://fly.io/blog/everyone-write-an-agent/) so you can internalize these concepts.

- context is precious
- degrades quickly after certain cliffs
- contexts are probably not going to get (much) bigger
- Why are you polluting it
- The context of the system is the most important thing you can put in your context window
- You, the engineer are the framer and arbiter of this

- How do we get people to write better prompts?
- How do we keep context windows clean?
- How do we generate better code?
- How do we continue to maintain the system in our heads?
- How do we create predictable factories of code that works, is flexible, maintainable, and secure?
- Stop dicking around with your IDE and code. Make no mistakes.

## It's Not The Size of the Token, but How You Use It

Many would tell you that's Claude's job now, the days of the engineer are waning and a sun is setting on the quirky guy who won't stop talking about Elixir. But that couldn't be farther from the truth for a very simple reason: context window limits.

Five years ago the engineer who directly transcribed customer diatribes into special cased code that solved that singular problem would have been anathema. Deeply distrusted by the rest of their team, their code consistently picked apart and asked to be rethought. And yet that has quickly become the expected way for us to interact with our LLM coding agents, spoon feeding them a couple sentences at a time, checking to make sure they're typing correctly, correcting tiny mistakes as we go, never reconsidering what was being asked. Would anyone sane treat their team of junior engineers so? Personally that sounds like an eternal hell I want no part of.

Currently humans are still necessary drivers of most software development.[^5] Contrary to popular belief this is not inherently due to the skill of humans, or any lack of skill or knowledge on behalf of LLMs. This is mainly due to the context size limitations of most models. Understanding what work needs to be done requires an astonishingly large amount of context, both in the traditinonal and token sense. How is a product or thing being used? Who is using it? What are they trying to accomplish? How are they failing to achieve their goals? What are the defeciencies of a system? Distilling these questions into coherent bodies of work that are properly prioritized requires a combination of very specific ways of thinking _as well as_ a wealth of experience. To a certain degree I'm of the opinion that

[^0]: I will say only this: I don't think there would be a similar reaction by companies if we still had sub 2% interest rates and high consumer confidence.

[^1]: I speak generally here, of course these things remain very important in domains like hardware or high throughput network routing. Most of us do not work in those domains.

[^2]: That is 19.2 minutes per PR assuming 8 hours of constant work... no breaks, no pauses, no discussion, nothing but pure prompt feeding. That leaves previous little time to scroll Twitter and post 10 times a day!

[^3]: I would love to be proven wrong.

[^4]: That is, as context grows the model becomes more efficacious _up to a very precise peak_, at which point things degrade rapidly.

[^0]: The syntax may be senior, but the brain is still very junior

[^1]: I'm not scarred, you're scarred!

[^2]: I mean "senior" in it's ancient sense, those that have weighty experience and lived to tell the tale to those who were not there

[^3]: I'm just as mad as the rest of you

[^4]: If that is possible!

[^5]: Feels a bit odd writing that
