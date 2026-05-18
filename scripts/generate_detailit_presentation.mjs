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
pptx.subject = 'Вводный курс для производственных подразделений';
pptx.title = 'ДетаЛит: безопасность, процесс, качество';
pptx.company = 'ООО "ДетаЛит"';
pptx.lang = 'ru-RU';
pptx.theme = {
  headFontFace: 'Arial',
  bodyFontFace: 'Arial',
  lang: 'ru-RU',
};
pptx.defineLayout({ name: 'WIDE', width: 13.333, height: 7.5 });

const colors = {
  navy: '1F2937',
  blue: '1F4E79',
  teal: '0F766E',
  amber: 'B45309',
  gray: '6B7280',
  pale: 'F5F7FA',
  white: 'FFFFFF',
};

function addTitle(slide, title, subtitle, accent = colors.blue) {
  slide.background = { color: colors.pale };
  slide.addShape(pptx.ShapeType.rect, { x: 0, y: 0, w: 13.333, h: 0.18, fill: { color: accent }, line: { color: accent } });
  slide.addText('ООО "ДетаЛит"', { x: 0.55, y: 0.38, w: 2.7, h: 0.28, fontFace: 'Arial', fontSize: 11, bold: true, color: accent, margin: 0 });
  slide.addText(title, { x: 0.55, y: 0.82, w: 11.8, h: 0.68, fontFace: 'Arial', fontSize: 29, bold: true, color: colors.navy, margin: 0 });
  slide.addText(subtitle, { x: 0.57, y: 1.55, w: 11.2, h: 0.34, fontFace: 'Arial', fontSize: 14, color: colors.gray, margin: 0 });
}

function addFooter(slide, index) {
  slide.addText(`Учебный материал • ${index}`, { x: 10.55, y: 7.12, w: 2.25, h: 0.2, fontSize: 8, color: colors.gray, align: 'right', margin: 0 });
}

function addCard(slide, x, y, w, h, title, body, accent) {
  slide.addShape(pptx.ShapeType.roundRect, {
    x, y, w, h,
    rectRadius: 0.08,
    fill: { color: colors.white },
    line: { color: 'DDE3EA', width: 1 },
  });
  slide.addShape(pptx.ShapeType.rect, { x, y, w: 0.08, h, fill: { color: accent }, line: { color: accent } });
  slide.addText(title, { x: x + 0.28, y: y + 0.22, w: w - 0.45, h: 0.32, fontSize: 15, bold: true, color: colors.navy, margin: 0 });
  slide.addText(body, { x: x + 0.28, y: y + 0.68, w: w - 0.45, h: h - 0.85, fontSize: 10.5, color: colors.gray, breakLine: false, fit: 'shrink', margin: 0.02 });
}

let slide = pptx.addSlide();
slide.background = { color: colors.navy };
slide.addShape(pptx.ShapeType.rect, { x: 0, y: 0, w: 13.333, h: 7.5, fill: { color: colors.navy }, line: { color: colors.navy } });
slide.addText('ООО "ДетаЛит"', { x: 0.7, y: 0.65, w: 5.8, h: 0.45, fontSize: 18, bold: true, color: 'A7F3D0', margin: 0 });
slide.addText('Безопасность, процесс и качество на производстве', { x: 0.7, y: 1.45, w: 9.8, h: 1.6, fontSize: 38, bold: true, color: colors.white, fit: 'shrink', margin: 0 });
slide.addText('Вводная презентация для сотрудников производственных подразделений', { x: 0.72, y: 3.16, w: 8.9, h: 0.46, fontSize: 16, color: 'D1D5DB', margin: 0 });
slide.addImage({ path: path.join(outDir, 'process-flow.png'), x: 7.2, y: 3.05, w: 5.25, h: 3.28 });
slide.addShape(pptx.ShapeType.rect, { x: 0.7, y: 6.6, w: 3.6, h: 0.08, fill: { color: 'A7F3D0' }, line: { color: 'A7F3D0' } });
addFooter(slide, 1);

slide = pptx.addSlide();
addTitle(slide, 'Карта рисков смены', 'Перед стартом сотрудник проверяет рабочее место, защиту и готовность оборудования.', colors.blue);
slide.addImage({ path: path.join(outDir, 'safety-zone.png'), x: 0.55, y: 2.15, w: 6.1, h: 3.82 });
addCard(slide, 7.05, 2.18, 5.55, 0.95, '1. Нет работы без СИЗ', 'Очки, перчатки, спецодежда и защитная обувь обязательны до входа на участок.', colors.blue);
addCard(slide, 7.05, 3.35, 5.55, 0.95, '2. Аварийные средства доступны', 'Проходы, кнопки стоп, аптечка и огнетушитель не перекрыты инструментом или тарой.', colors.blue);
addCard(slide, 7.05, 4.52, 5.55, 0.95, '3. Отклонение фиксируется', 'Неисправность, перегрев, шум, вибрация и запах гари сразу передаются мастеру.', colors.blue);
addFooter(slide, 2);

slide = pptx.addSlide();
addTitle(slide, 'Производственный маршрут', 'Каждая партия должна проходить один и тот же контролируемый путь.', colors.teal);
const route = [
  ['Задание смены', 'Получить маршрутную карту, чертеж, материал и план партии.'],
  ['Форма', 'Проверить чистоту формы, смазку, крепление и маркировку.'],
  ['Плавка', 'Сверить температуру и состав с технологической картой.'],
  ['Заливка', 'Вести процесс стабильно, без рывков и обхода защит.'],
  ['ОТК', 'Проверить дефекты, геометрию, маркировку и решение по партии.'],
];
route.forEach((item, i) => {
  const x = 0.75 + i * 2.5;
  slide.addShape(pptx.ShapeType.chevron, { x, y: 2.55, w: 2.2, h: 1.05, fill: { color: i % 2 === 0 ? colors.teal : '115E59' }, line: { color: colors.teal } });
  slide.addText(String(i + 1), { x: x + 0.14, y: 2.88, w: 0.3, h: 0.28, fontSize: 12, bold: true, color: colors.white, align: 'center', margin: 0 });
  slide.addText(item[0], { x: x + 0.45, y: 2.8, w: 1.42, h: 0.28, fontSize: 12, bold: true, color: colors.white, align: 'center', margin: 0 });
  slide.addText(item[1], { x: x - 0.03, y: 3.88, w: 2.22, h: 0.75, fontSize: 9.5, color: colors.gray, align: 'center', margin: 0.02, fit: 'shrink' });
});
slide.addText('Ключевая мысль: если один этап пропущен, качество партии нельзя доказать документально.', { x: 1.0, y: 5.55, w: 11.1, h: 0.45, fontSize: 17, bold: true, color: colors.navy, align: 'center', margin: 0 });
addFooter(slide, 3);

slide = pptx.addSlide();
addTitle(slide, 'Как отличать годное изделие от брака', 'Контроль качества строится на видимых признаках, измерениях и прослеживаемости.', colors.amber);
slide.addImage({ path: path.join(outDir, 'quality-checkpoints.png'), x: 6.82, y: 2.12, w: 5.85, h: 3.65 });
addCard(slide, 0.7, 2.05, 5.55, 0.9, 'Поверхность', 'Нет трещин, сквозных раковин, непроливов, подгара и следов перегрева.', colors.amber);
addCard(slide, 0.7, 3.15, 5.55, 0.9, 'Размеры', 'Критические размеры сверяются по чертежу, результат фиксируется в журнале ОТК.', colors.amber);
addCard(slide, 0.7, 4.25, 5.55, 0.9, 'Маркировка', 'Изделие должно быть связано с партией, сменой, маршрутной картой и ответственным.', colors.amber);
addCard(slide, 0.7, 5.35, 5.55, 0.9, 'Спорные случаи', 'Сотрудник не принимает спорное изделие сам: партия отделяется, решение принимает мастер и ОТК.', colors.amber);
addFooter(slide, 4);

slide = pptx.addSlide();
addTitle(slide, 'Что запомнить после курса', 'Минимальный стандарт поведения сотрудника на участке.', colors.blue);
[
  ['Остановить опасную работу', 'Если есть риск травмы, перегрева, пожара или повреждения оборудования.'],
  ['Не скрывать отклонения', 'Своевременная фиксация дешевле переделки партии и спора с заказчиком.'],
  ['Проверять документы', 'Чертеж, маршрутная карта и маркировка партии так же важны, как сама деталь.'],
  ['Обращаться к мастеру', 'Любой спорный случай решается через мастера смены и ОТК.'],
].forEach((item, i) => {
  const x = 0.75 + (i % 2) * 6.05;
  const y = 2.2 + Math.floor(i / 2) * 1.7;
  addCard(slide, x, y, 5.5, 1.1, item[0], item[1], i % 2 === 0 ? colors.blue : colors.teal);
});
slide.addText('Финальная проверка в курсе закрепляет эти правила на практических ситуациях.', { x: 0.8, y: 6.1, w: 11.8, h: 0.34, fontSize: 14, color: colors.gray, align: 'center', margin: 0 });
addFooter(slide, 5);

await pptx.writeFile({ fileName: path.join(outDir, 'detailit-production-intro.pptx') });
console.log(`Generated ${path.join(outDir, 'detailit-production-intro.pptx')}`);
