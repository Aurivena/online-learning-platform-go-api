import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { createRequire } from 'node:module';

const require = createRequire(import.meta.url);
const pptxgen = require('pptxgenjs');

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(__dirname, '..');
const outDir = path.join(root, 'resources', 'course_files', 'detailit');

const pptx = new pptxgen();
pptx.layout = 'LAYOUT_WIDE';
pptx.author = 'ООО "ДетаЛит"';
pptx.subject = 'Курсы по ОКВЭД 29.31';
pptx.title = 'ДетаЛит: автоэлектроника, жгуты, ESD и испытания';
pptx.company = 'ООО "ДетаЛит"';
pptx.lang = 'ru-RU';
pptx.theme = { headFontFace: 'Arial', bodyFontFace: 'Arial', lang: 'ru-RU' };

const c = {
  ink: '1F2937',
  muted: '64748B',
  pale: 'F8FAFC',
  white: 'FFFFFF',
  cyan: '155E75',
  teal: '0F766E',
  violet: '7C3AED',
  amber: 'B45309',
};

function title(slide, text, sub, accent) {
  slide.background = { color: c.pale };
  slide.addShape(pptx.ShapeType.rect, { x: 0, y: 0, w: 13.333, h: 0.18, fill: { color: accent }, line: { color: accent } });
  slide.addText('ООО "ДетаЛит"', { x: 0.55, y: 0.36, w: 3.1, h: 0.25, fontSize: 11, bold: true, color: accent, margin: 0 });
  slide.addText(text, { x: 0.55, y: 0.82, w: 11.9, h: 0.72, fontSize: 29, bold: true, color: c.ink, fit: 'shrink', margin: 0 });
  slide.addText(sub, { x: 0.56, y: 1.55, w: 11.6, h: 0.34, fontSize: 13.5, color: c.muted, margin: 0 });
}

function card(slide, x, y, w, h, head, body, accent) {
  slide.addShape(pptx.ShapeType.roundRect, { x, y, w, h, rectRadius: 0.08, fill: { color: c.white }, line: { color: 'D8E0EA', width: 1 } });
  slide.addShape(pptx.ShapeType.rect, { x, y, w: 0.08, h, fill: { color: accent }, line: { color: accent } });
  slide.addText(head, { x: x + 0.25, y: y + 0.18, w: w - 0.45, h: 0.3, fontSize: 14, bold: true, color: c.ink, margin: 0 });
  slide.addText(body, { x: x + 0.25, y: y + 0.58, w: w - 0.45, h: h - 0.68, fontSize: 10, color: c.muted, fit: 'shrink', margin: 0.02 });
}

function footer(slide, n) {
  slide.addText(`Учебная сетка 29.31 • ${n}`, { x: 10.4, y: 7.12, w: 2.35, h: 0.2, fontSize: 8, color: c.muted, align: 'right', margin: 0 });
}

let s = pptx.addSlide();
s.background = { color: c.ink };
s.addText('ООО "ДетаЛит"', { x: 0.7, y: 0.68, w: 4.3, h: 0.45, fontSize: 18, bold: true, color: 'A7F3D0', margin: 0 });
s.addText('Электронное оборудование для автотранспортных средств', { x: 0.7, y: 1.4, w: 10.1, h: 1.45, fontSize: 35, bold: true, color: c.white, fit: 'shrink', margin: 0 });
s.addText('Учебная презентация для подразделений: производство, жгуты, ESD, испытания и ОТК', { x: 0.72, y: 3.0, w: 8.8, h: 0.5, fontSize: 15.5, color: 'CBD5E1', margin: 0 });
s.addImage({ path: path.join(outDir, 'electronics-architecture.png'), x: 6.9, y: 3.15, w: 5.55, h: 3.47 });
s.addShape(pptx.ShapeType.rect, { x: 0.7, y: 6.58, w: 3.9, h: 0.08, fill: { color: 'A7F3D0' }, line: { color: 'A7F3D0' } });
footer(s, 1);

s = pptx.addSlide();
title(s, 'Как устроена учебная сетка', 'Курсы связаны с подразделениями и рабочими ролями, а не висят отдельными карточками.', c.cyan);
[
  ['Производство электрооборудования', 'Базовая архитектура изделия, маршрут партии, критичные операции и допуск к работе.'],
  ['Сборка жгутов и разъемов', 'Провод, контакт, обжим, пиновка, маркировка, финальная проверка цепей.'],
  ['ESD и монтаж', 'Защита компонентов, пайка, визуальный контроль, работа с платами и скрытыми отказами.'],
  ['Испытания и трассируемость', 'Методики измерений, запись результата, несоответствия и история партии.'],
].forEach((item, i) => card(s, 0.75 + (i % 2) * 6.05, 2.25 + Math.floor(i / 2) * 1.62, 5.55, 1.05, item[0], item[1], [c.cyan, c.teal, c.violet, c.amber][i]));
footer(s, 2);

s = pptx.addSlide();
title(s, 'Автоэлектроника: от сигнала к изделию', 'Сотрудник должен понимать, почему качество маленькой операции влияет на надежность автомобиля.', c.cyan);
s.addImage({ path: path.join(outDir, 'electronics-architecture.png'), x: 0.65, y: 2.1, w: 6.25, h: 3.9 });
card(s, 7.22, 2.2, 5.25, 0.9, 'Сигнал', 'Датчик или цепь должны передать корректное состояние без шумов, обрывов и переполюсовки.', c.cyan);
card(s, 7.22, 3.32, 5.25, 0.9, 'Управление', 'ЭБУ выполняет алгоритм, а исполнительная цепь должна сработать в заданном режиме.', c.cyan);
card(s, 7.22, 4.44, 5.25, 0.9, 'Надежность', 'Ошибка в жгуте, пайке или маркировке превращается в отказ уже у заказчика.', c.cyan);
footer(s, 3);

s = pptx.addSlide();
title(s, 'Жгуты и разъемы', 'Проверка нужна на каждом шаге: после обжима исправлять ошибку гораздо дороже.', c.teal);
s.addImage({ path: path.join(outDir, 'harness-assembly.png'), x: 6.72, y: 2.08, w: 5.85, h: 3.65 });
[
  ['Подготовка', 'Сечение, длина, цвет, маркировка и целостность изоляции сверяются до операции.'],
  ['Обжим', 'Контакт без повреждения жилы, с правильной геометрией и усилием удержания.'],
  ['Пиновка', 'Провод стоит в правильной ячейке разъема; ошибка блокирует выпуск партии.'],
  ['Тест цепи', 'Непрерывность, отсутствие короткого замыкания и запись результата.'],
].forEach((item, i) => card(s, 0.7, 2.05 + i * 1.08, 5.45, 0.82, item[0], item[1], c.teal));
footer(s, 4);

s = pptx.addSlide();
title(s, 'ESD, пайка и монтаж', 'Скрытые повреждения компонентов часто не видны сразу, поэтому пост должен быть дисциплинированным.', c.violet);
s.addImage({ path: path.join(outDir, 'esd-workstation.png'), x: 0.75, y: 2.05, w: 5.95, h: 3.72 });
card(s, 7.0, 2.05, 5.55, 0.95, 'ESD до операции', 'Проверить браслет, коврик, тару, заземление и отсутствие обычного пластика на посту.', c.violet);
card(s, 7.0, 3.25, 5.55, 0.95, 'Пайка по карте', 'Температура, флюс, время прогрева и визуальные критерии берутся из техкарты.', c.violet);
card(s, 7.0, 4.45, 5.55, 0.95, 'Спорная пайка', 'Не передается дальше: повторный осмотр, доработка или оформление несоответствия.', c.violet);
footer(s, 5);

s = pptx.addSlide();
title(s, 'Испытания и трассируемость', 'Качество надо не только сделать, но и доказать по данным партии.', c.amber);
s.addImage({ path: path.join(outDir, 'testing-traceability.png'), x: 6.75, y: 2.08, w: 5.85, h: 3.65 });
card(s, 0.75, 2.15, 5.35, 0.9, 'Методика', 'Каждое измерение имеет допуск, прибор, порядок выполнения и ответственного.', c.amber);
card(s, 0.75, 3.3, 5.35, 0.9, 'Запись', 'Номер партии связывает изделие с материалом, операцией, сменой и сотрудником.', c.amber);
card(s, 0.75, 4.45, 5.35, 0.9, 'Несоответствие', 'Изделие изолируется, причина фиксируется, решение принимает ОТК и мастер.', c.amber);
footer(s, 6);

await pptx.writeFile({ fileName: path.join(outDir, 'detailit-auto-electronics-training.pptx') });
console.log(`Generated ${path.join(outDir, 'detailit-auto-electronics-training.pptx')}`);
